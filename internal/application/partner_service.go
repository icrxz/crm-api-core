package application

import (
	"context"
	"io"
	"strings"

	"github.com/icrxz/crm-api-core/internal/domain"
)

type partnerService struct {
	partnerRepository domain.PartnerRepository
}

type PartnerService interface {
	Create(ctx context.Context, partner domain.Partner) (string, error)
	GetByID(ctx context.Context, partnerID string) (*domain.Partner, error)
	Update(ctx context.Context, partnerID string, editPartner domain.EditPartner) error
	Delete(ctx context.Context, partnerID string) error
	Search(ctx context.Context, filters domain.PartnerFilters) (domain.PagingResult[domain.Partner], error)
	CreateBatch(ctx context.Context, file io.Reader, createdBy string) ([]string, error)
}

func NewPartnerService(partnerRepository domain.PartnerRepository) PartnerService {
	return &partnerService{
		partnerRepository: partnerRepository,
	}
}

func (s *partnerService) Create(ctx context.Context, partner domain.Partner) (string, error) {
	return s.partnerRepository.Create(ctx, partner)
}

func (s *partnerService) GetByID(ctx context.Context, partnerID string) (*domain.Partner, error) {
	if partnerID == "" {
		return nil, domain.NewValidationError("partnerID cannot be empty", nil)
	}

	return s.partnerRepository.GetByID(ctx, partnerID)
}

func (s *partnerService) Update(ctx context.Context, partnerID string, editPartner domain.EditPartner) error {
	if partnerID == "" {
		return domain.NewValidationError("partnerID cannot be empty", nil)
	}

	partner, err := s.partnerRepository.GetByID(ctx, partnerID)
	if err != nil {
		return err
	}

	partner.MergeUpdate(editPartner)

	return s.partnerRepository.Update(ctx, *partner)
}

func (s *partnerService) Delete(ctx context.Context, partnerID string) error {
	if partnerID == "" {
		return domain.NewValidationError("partnerID cannot be empty", nil)
	}

	return s.partnerRepository.Delete(ctx, partnerID)
}

func (s *partnerService) Search(ctx context.Context, filters domain.PartnerFilters) (domain.PagingResult[domain.Partner], error) {
	return s.partnerRepository.Search(ctx, filters)
}

func (s *partnerService) CreateBatch(ctx context.Context, file io.Reader, createdBy string) ([]string, error) {
	partnersRows, err := readCSV(file)
	if err != nil {
		return nil, err
	}

	columnsIndex := getColumnHeadersIndex(partnersRows[0])
	partners, err := s.buildPartner(partnersRows[1:], columnsIndex, createdBy)
	if err != nil {
		return nil, err
	}

	return s.partnerRepository.CreateBatch(ctx, partners)
}

func (s *partnerService) buildPartner(csvRows [][]string, columnsIndex map[string]int, author string) ([]domain.Partner, error) {
	partners := make([]domain.Partner, 0, len(csvRows))

	for _, row := range csvRows {
		phone := row[columnsIndex["Telefone"]]
		if strings.TrimSpace(phone) != "" {
			phone = "+55 " + phone
		}

		personalContact := domain.Contact{
			PhoneNumber: phone,
		}

		shippingAddress := domain.Address{
			City:    row[columnsIndex["Cidade"]],
			State:   row[columnsIndex["Estado"]],
			Country: "brazil",
		}

		document := row[columnsIndex["Documento"]]
		document = strings.ReplaceAll(document, ".", "")
		document = strings.ReplaceAll(document, "-", "")
		document = strings.ReplaceAll(document, "/", "")

		documentType := ""
		if len(document) == 11 {
			documentType = "CPF"
		} else if len(document) == 14 {
			documentType = "CNPJ"
		}

		description := row[columnsIndex["Observacoes"]]

		newPartner, err := domain.NewPartner(
			row[columnsIndex["Nome"]],
			row[columnsIndex["Sobrenome"]],
			"",
			"",
			document,
			documentType,
			author,
			personalContact,
			domain.Contact{},
			shippingAddress,
			domain.Address{},
			description,
			row[columnsIndex["Tipo"]],
			"",
			"",
		)
		if err != nil {
			return nil, err
		}

		partners = append(partners, newPartner)
	}
	return partners, nil
}
