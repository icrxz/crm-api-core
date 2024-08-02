package application

import (
	"context"
	"encoding/csv"
	"github.com/icrxz/crm-api-core/internal/domain"
	"io"
)

type partnerService struct {
	partnerRepository domain.PartnerRepository
}

type PartnerService interface {
	Create(ctx context.Context, partner domain.Partner) (string, error)
	GetByID(ctx context.Context, partnerID string) (*domain.Partner, error)
	Update(ctx context.Context, partnerID string, editPartner domain.EditPartner) error
	Delete(ctx context.Context, partnerID string) error
	Search(ctx context.Context, filters domain.PartnerFilters) ([]domain.Partner, error)
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

func (s *partnerService) Search(ctx context.Context, filters domain.PartnerFilters) ([]domain.Partner, error) {
	return s.partnerRepository.Search(ctx, filters)
}

func (s *partnerService) CreateBatch(ctx context.Context, file io.Reader, createdBy string) ([]string, error) {
	fileCSV := csv.NewReader(file)

	partnersRows, err := s.readCSV(fileCSV)
	if err != nil {
		return nil, err
	}

	columnsIndex := s.getColumnHeadersIndex(partnersRows[0])
	partners, err := s.buildPartner(partnersRows[1:], columnsIndex, createdBy)
	if err != nil {
		return nil, err
	}

	return s.partnerRepository.CreateBatch(ctx, partners)
}

func (s *partnerService) readCSV(fileCSV *csv.Reader) ([][]string, error) {
	csvRows := make([][]string, 0)

	for {
		row, err := fileCSV.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		csvRows = append(csvRows, row)
	}

	return csvRows, nil
}

func (s *partnerService) buildPartner(csvRows [][]string, columnsIndex map[string]int, author string) ([]domain.Partner, error) {
	partners := make([]domain.Partner, 0, len(csvRows))

	for _, row := range csvRows {
		personalContact := domain.Contact{
			PhoneNumber: row[columnsIndex["telefone"]],
			Email:       row[columnsIndex["email"]],
		}

		shippingAddress := domain.Address{
			City:    row[columnsIndex["cidade"]],
			State:   row[columnsIndex["estado"]],
			Country: "Brasil",
		}

		newPartner, err := domain.NewPartner(
			row[columnsIndex["nome"]],
			row[columnsIndex["sobrenome"]],
			"",
			"",
			row[columnsIndex["documento"]],
			"",
			author,
			personalContact,
			domain.Contact{},
			shippingAddress,
			domain.Address{},
			row[columnsIndex["observacoes"]],
		)
		if err != nil {
			return nil, err
		}

		partners = append(partners, newPartner)
	}
	return partners, nil
}

func (s *partnerService) getColumnHeadersIndex(header []string) map[string]int {
	columnsIndex := make(map[string]int)
	for i, column := range header {
		columnsIndex[column] = i
	}
	return columnsIndex
}
