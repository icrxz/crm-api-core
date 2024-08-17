package application

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/icrxz/crm-api-core/internal/domain"
	"github.com/nguyenthenguyen/docx"
	"golang.org/x/sync/errgroup"
)

const (
	dateReportLayout     = "02/Jan/2006"
	dateTimeReportLayout = "02/Jan/2006 15:04"
	timestampLayout      = "02_01_2006_15_04_05_0000"
)

var contractorsTemplates = map[string]string{
	"LuizaSeg":     "luizaseg_template.docx",
	"Assurant":     "assurant_template.docx",
	"Cardif":       "cardif_template.docx",
	"Ezze Seguros": "ezze_template.docx",
}

type ContentWithAttachment struct {
	Content    string
	Attachment [][]byte
}

type reportService struct {
	reportFolder      string
	caseService       CaseService
	productService    ProductService
	customerService   CustomerService
	commentService    CommentService
	partnerService    PartnerService
	contractorService ContractorService
	attachmentBucket  domain.AttachmentBucket
}

type ReportService interface {
	GenerateReport(ctx context.Context, crmCase domain.Case) ([]byte, string, error)
}

type ReportData struct {
	CrmCase    domain.Case
	Customer   domain.Customer
	Product    domain.Product
	Partner    domain.Partner
	Contractor domain.Contractor
	Comments   []domain.Comment
}

func NewReportService(
	reportFolder string,
	caseService CaseService,
	productService ProductService,
	customerService CustomerService,
	commentService CommentService,
	partnerService PartnerService,
	contractorService ContractorService,
	attachmentBucket domain.AttachmentBucket,
) ReportService {
	return &reportService{
		reportFolder:      reportFolder,
		caseService:       caseService,
		productService:    productService,
		customerService:   customerService,
		commentService:    commentService,
		partnerService:    partnerService,
		contractorService: contractorService,
		attachmentBucket:  attachmentBucket,
	}
}

func (s *reportService) GenerateReport(ctx context.Context, crmCase domain.Case) ([]byte, string, error) {
	var memoryDoc bytes.Buffer

	reportData, err := s.getReportData(ctx, crmCase)
	if err != nil {
		return nil, "", err
	}

	filename, hasTemplate := contractorsTemplates[reportData.Contractor.CompanyName]
	if !hasTemplate {
		return nil, "", fmt.Errorf("no template found for contractor %s", reportData.Contractor.CompanyName)
	}

	err = s.readReportTemplate(ctx, *reportData, filename, &memoryDoc)
	if err != nil {
		return nil, "", err
	}

	return memoryDoc.Bytes(), fmt.Sprintf("%s-%s-%s", reportData.Contractor.CompanyName, crmCase.ExternalReference, time.Now().Format(timestampLayout)), nil
}

func (s *reportService) getReportData(ctx context.Context, crmCase domain.Case) (*ReportData, error) {
	reportData := &ReportData{
		CrmCase: crmCase,
	}

	wg, newCtx := errgroup.WithContext(ctx)
	wg.Go(func() error {
		customer, err := s.customerService.GetByID(newCtx, crmCase.CustomerID)
		if err != nil {
			return err
		}
		reportData.Customer = *customer
		return nil
	})

	wg.Go(func() error {
		product, err := s.productService.GetProductByID(newCtx, crmCase.ProductID)
		if err != nil {
			return err
		}
		reportData.Product = *product
		return nil
	})

	wg.Go(func() error {
		comments, err := s.commentService.GetByCaseID(newCtx, crmCase.CaseID)
		if err != nil {
			return err
		}
		reportData.Comments = comments
		return nil
	})

	wg.Go(func() error {
		partner, err := s.partnerService.GetByID(newCtx, crmCase.PartnerID)
		if err != nil {
			return err
		}
		reportData.Partner = *partner
		return nil
	})

	wg.Go(func() error {
		contractor, err := s.contractorService.GetByID(newCtx, crmCase.ContractorID)
		if err != nil {
			return err
		}
		reportData.Contractor = *contractor
		return nil
	})

	if err := wg.Wait(); err != nil {
		return nil, err
	}

	return reportData, nil
}

func (s *reportService) readReportTemplate(ctx context.Context, reportData ReportData, filename string, memDoc io.Writer) error {
	filePath := fmt.Sprintf("%s/%s", s.reportFolder, filename)
	file, err := docx.ReadDocxFile(filePath)
	if err != nil {
		return err
	}

	docEdit := file.Editable()
	defer file.Close()

	err = docEdit.Replace("$claim", reportData.CrmCase.ExternalReference, -1)
	if err != nil {
		return err
	}

	err = docEdit.Replace("$actual_date", time.Now().Format(dateReportLayout), -1)
	if err != nil {
		return err
	}

	err = docEdit.Replace("$client", fmt.Sprintf("%s %s", reportData.Customer.FirstName, reportData.Customer.LastName), -1)
	if err != nil {
		return err
	}

	err = docEdit.Replace("$brand", reportData.Product.Brand, -1)
	if err != nil {
		return err
	}

	err = docEdit.Replace("$summary", reportData.CrmCase.Subject, -1)
	if err != nil {
		return err
	}
	err = docEdit.Replace("$partner", fmt.Sprintf("%s %s", reportData.Partner.FirstName, reportData.Partner.LastName), -1)
	if err != nil {
		return err
	}
	err = docEdit.Replace("$target_date", reportData.CrmCase.TargetDate.Format(dateReportLayout), -1)
	if err != nil {
		return err
	}
	err = docEdit.Replace("$document", ParseDocument(reportData.Customer.Document), -1)
	if err != nil {
		return err
	}
	err = docEdit.Replace("$address", reportData.Customer.ShippingAddress.Address, -1)
	if err != nil {
		return err
	}
	err = docEdit.Replace("$zip_code", reportData.Customer.ShippingAddress.ZipCode, -1)
	if err != nil {
		return err
	}
	err = docEdit.Replace("$product", reportData.Product.Name, -1)
	if err != nil {
		return err
	}
	err = docEdit.Replace("$serial_number", reportData.Product.SerialNumber, -1)
	if err != nil {
		return err
	}

	content := make([]string, 0)
	contentAttachments := make([][]byte, 0)
	comments := make([]string, 0)
	commentsAttachments := make([][]byte, 0)
	resolution := make([]string, 0)
	resolutionAttachments := make([][]byte, 0)
	for _, comment := range reportData.Comments {
		switch comment.CommentType {
		case domain.CONTENT:
			contentAttachments, err = s.downloadFiles(ctx, comment.Attachments)
			if err != nil {
				return err
			}
			content = append(content, fmt.Sprintf("%s - %s", comment.CreatedAt.Format(dateTimeReportLayout), comment.Content))
		case domain.RESOLUTION:
			resolutionAttachments, err = s.downloadFiles(ctx, comment.Attachments)
			if err != nil {
				return err
			}
			resolution = append(resolution, fmt.Sprintf("%s - %s", comment.CreatedAt.Format(dateTimeReportLayout), comment.Content))
		case domain.COMMENT:
			commentsAttachments, err = s.downloadFiles(ctx, comment.Attachments)
			if err != nil {
				return err
			}
			comments = append(comments, fmt.Sprintf("%s - %s", comment.CreatedAt.Format(dateTimeReportLayout), comment.Content))
		}
	}

	err = docEdit.Replace("$content", strings.Join(content, "\r\n"), -1)
	if err != nil {
		return err
	}
	if len(contentAttachments) > 0 {
		err = docEdit.Replace("$image_content", string(contentAttachments[0]), -1)
		if err != nil {
			return err
		}
	}

	err = docEdit.Replace("$comments", strings.Join(comments, "\r\n"), -1)
	if err != nil {
		return err
	}
	if len(commentsAttachments) > 0 {
		err = docEdit.Replace("$image_comment", string(commentsAttachments[0]), -1)
		if err != nil {
			return err
		}
	}

	err = docEdit.Replace("$resolution", strings.Join(resolution, "\r\n"), -1)
	if err != nil {
		return err
	}
	if len(resolutionAttachments) > 0 {
		err = docEdit.Replace("$image_resolution", string(resolutionAttachments[0]), -1)
		if err != nil {
			return err
		}
	}

	return docEdit.Write(memDoc)
}

func (s *reportService) downloadFiles(ctx context.Context, files []domain.Attachment) ([][]byte, error) {
	downloadedFiles := make([][]byte, 0)
	for _, attachment := range files {
		file, err := s.attachmentBucket.Download(ctx, attachment.Key)
		if err != nil {
			return nil, err
		}
		downloadedFiles = append(downloadedFiles, file)
	}

	return downloadedFiles, nil
}
