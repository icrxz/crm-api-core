package application

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	_ "image/png"
	"io"
	"os"
	"time"

	"github.com/google/uuid"
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

	err = docEdit.Replace("$model", reportData.Product.Model, -1)
	if err != nil {
		return err
	}

	var resolution string
	resolutionAttachments := make([][]byte, 0)

	for _, comment := range reportData.Comments {
		switch comment.CommentType {
		case domain.COMMENT_RESOLUTION:
			resolutionAttachments, err = s.downloadFiles(ctx, comment.Attachments)
			if err != nil {
				return err
			}
		case domain.COMMENT_REPORT:
			resolution = comment.Content
		}
	}

	err = docEdit.Replace("$resolution", fmt.Sprintf("%s\r\n", resolution), -1)
	if err != nil {
		return err
	}

	isAssurant := reportData.Contractor.CompanyName == "Assurant"
	err = s.replaceImages(docEdit, memDoc, resolutionAttachments, isAssurant)
	if err != nil {
		return err
	}

	return nil
}

func (s *reportService) downloadFiles(ctx context.Context, files []domain.Attachment) ([][]byte, error) {
	downloadedFiles := make([][]byte, 0)
	for _, attachment := range files {
		fmt.Println("")
		file, err := s.attachmentBucket.Download(ctx, attachment.Key)
		if err != nil {
			return nil, err
		}
		downloadedFiles = append(downloadedFiles, file)
	}

	return downloadedFiles, nil
}

func (s *reportService) replaceImages(doc *docx.Docx, memDoc io.Writer, attachments [][]byte, isAssurant bool) error {
	attachmentNames := make([]string, 0, len(attachments))

	docImagesLength := doc.ImagesLen()

	for idx := range docImagesLength {
		if idx >= docImagesLength-1 && isAssurant {
			break
		}

		if idx >= len(attachments) {
			break
		}

		attachment := attachments[idx]

		img, _, err := image.Decode(bytes.NewReader(attachment))
		if err != nil {
			return err
		}

		fileName := fmt.Sprintf("./resources/img_%s.png", uuid.NewString())
		out, _ := os.Create(fileName)
		attachmentNames = append(attachmentNames, fileName)

		err = jpeg.Encode(out, img, nil)
		if err != nil {
			return err
		}

		err = out.Close()
		if err != nil {
			return err
		}

		err = doc.ReplaceImage(fmt.Sprintf("word/media/image%d.png", idx+1), fileName)
		if err != nil {
			return err
		}
	}

	err := doc.Write(memDoc)
	if err != nil {
		return err
	}

	for _, attachmentName := range attachmentNames {
		err = os.Remove(attachmentName)
		if err != nil {
			return err
		}
	}

	return nil
}
