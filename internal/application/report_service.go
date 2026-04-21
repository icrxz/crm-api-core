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

	_ "golang.org/x/image/webp"

	"github.com/google/uuid"
	"github.com/icrxz/crm-api-core/internal/domain"
	"github.com/nguyenthenguyen/docx"
	"golang.org/x/sync/errgroup"
)

const (
	dateReportLayout = "02/Jan/2006"
	timestampLayout  = "02_01_2006_15_04_05_0000"
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

//go:generate mockgen -source=report_service.go -destination=mock_application/mock_report_service.go -package=mock_application
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
		return fmt.Errorf("failed to read report template %s: %w", filename, err)
	}

	docEdit := file.Editable()
	defer file.Close()

	if err = s.replaceReportFields(docEdit, reportData); err != nil {
		return err
	}

	resolution, resolutionAttachments, attachmentKeys, err := s.extractResolutionFromComments(ctx, reportData.Comments)
	if err != nil {
		return err
	}

	if err = docEdit.Replace("$resolution", fmt.Sprintf("%s\r\n", resolution), -1); err != nil {
		return err
	}

	isAssurant := reportData.Contractor.CompanyName == "Assurant"
	return s.replaceImages(docEdit, memDoc, resolutionAttachments, attachmentKeys, isAssurant)
}

func (s *reportService) replaceReportFields(docEdit *docx.Docx, reportData ReportData) error {
	replacements := []struct {
		placeholder string
		value       string
	}{
		{"$claim", reportData.CrmCase.ExternalReference},
		{"$actual_date", time.Now().Format(dateReportLayout)},
		{"$client", fmt.Sprintf("%s %s", reportData.Customer.FirstName, reportData.Customer.LastName)},
		{"$brand", reportData.Product.Brand},
		{"$summary", reportData.CrmCase.Subject},
		{"$partner", fmt.Sprintf("%s %s", reportData.Partner.FirstName, reportData.Partner.LastName)},
		{"$target_date", reportData.CrmCase.TargetDate.Format(dateReportLayout)},
		{"$document", ParseDocument(reportData.Customer.Document)},
		{"$address", reportData.Customer.ShippingAddress.Address},
		{"$zip_code", reportData.Customer.ShippingAddress.ZipCode},
		{"$product", reportData.Product.Name},
		{"$serial_number", reportData.Product.SerialNumber},
		{"$model", reportData.Product.Model},
	}

	for _, r := range replacements {
		if err := docEdit.Replace(r.placeholder, r.value, -1); err != nil {
			return fmt.Errorf("failed to replace field %s: %w", r.placeholder, err)
		}
	}

	return nil
}

func (s *reportService) extractResolutionFromComments(ctx context.Context, comments []domain.Comment) (string, [][]byte, []string, error) {
	var resolution string
	resolutionAttachments := make([][]byte, 0)
	attachmentKeys := make([]string, 0)

	for _, comment := range comments {
		switch comment.CommentType {
		case domain.COMMENT_RESOLUTION:
			files, err := s.downloadFiles(ctx, comment.Attachments)
			if err != nil {
				return "", nil, nil, err
			}
			resolutionAttachments = files
			for _, a := range comment.Attachments {
				attachmentKeys = append(attachmentKeys, a.Key)
			}
		case domain.COMMENT_REPORT:
			resolution = comment.Content
		}
	}

	return resolution, resolutionAttachments, attachmentKeys, nil
}

func (s *reportService) downloadFiles(ctx context.Context, files []domain.Attachment) ([][]byte, error) {
	downloadedFiles := make([][]byte, 0)
	for _, attachment := range files {
		file, err := s.attachmentBucket.Download(ctx, attachment.Key)
		if err != nil {
			return nil, fmt.Errorf("failed to download attachment %s: %w", attachment.Key, err)
		}
		downloadedFiles = append(downloadedFiles, file)
	}

	return downloadedFiles, nil
}

func (s *reportService) replaceImages(doc *docx.Docx, memDoc io.Writer, attachments [][]byte, keys []string, isAssurant bool) error {
	attachmentNames := make([]string, 0, len(attachments))

	docImagesLength := doc.ImagesLen()

	for idx := range docImagesLength {
		if idx >= docImagesLength-1 && isAssurant {
			break
		}

		if idx >= len(attachments) {
			break
		}

		key := keys[idx]
		img, _, err := image.Decode(bytes.NewReader(attachments[idx]))
		if err != nil {
			return fmt.Errorf("failed to decode attachment %s: %w", key, err)
		}

		fileName := fmt.Sprintf("./resources/img_%s.png", uuid.NewString())
		out, err := os.Create(fileName)
		if err != nil {
			return fmt.Errorf("failed to create temp file for attachment %s: %w", key, err)
		}
		attachmentNames = append(attachmentNames, fileName)

		if err = jpeg.Encode(out, img, nil); err != nil {
			return fmt.Errorf("failed to encode attachment %s as jpeg: %w", key, err)
		}

		if err = out.Close(); err != nil {
			return fmt.Errorf("failed to close temp file for attachment %s: %w", key, err)
		}

		if err = doc.ReplaceImage(fmt.Sprintf("word/media/image%d.png", idx+1), fileName); err != nil {
			return fmt.Errorf("failed to replace image slot %d with attachment %s: %w", idx+1, key, err)
		}
	}

	if err := doc.Write(memDoc); err != nil {
		return fmt.Errorf("failed to write report document: %w", err)
	}

	for _, attachmentName := range attachmentNames {
		if err := os.Remove(attachmentName); err != nil {
			return fmt.Errorf("failed to remove temp file %s: %w", attachmentName, err)
		}
	}

	return nil
}
