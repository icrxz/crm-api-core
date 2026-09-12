package application

import (
	"context"
	"testing"

	"github.com/icrxz/crm-api-core/internal/domain"
	"github.com/icrxz/crm-api-core/internal/domain/mock_domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestPartnerService_Create(t *testing.T) {
	t.Run("returns validation error when document is empty", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		service := NewPartnerService(mock_domain.NewMockPartnerRepository(ctrl))

		partner := domain.Partner{DocumentType: domain.CPF}

		_, err := service.Create(context.Background(), partner)

		require.Error(t, err)
	})

	t.Run("returns validation error when document_type is not CPF or CNPJ", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		service := NewPartnerService(mock_domain.NewMockPartnerRepository(ctrl))

		partner := domain.Partner{Document: "12345678900", DocumentType: domain.RG}

		_, err := service.Create(context.Background(), partner)

		require.Error(t, err)
	})

	t.Run("returns conflict when the document already belongs to another partner", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mock_domain.NewMockPartnerRepository(ctrl)
		service := NewPartnerService(mockRepo)

		partner := domain.Partner{Document: "12345678900", DocumentType: domain.CPF}

		mockRepo.EXPECT().Search(gomock.Any(), domain.PartnerFilters{
			DocumentExact: []string{"12345678900"},
			PagingFilter:  domain.PagingFilter{Limit: 1, Offset: 0},
		}).Return(domain.PagingResult[domain.Partner]{
			Result: []domain.Partner{{PartnerID: "partner-1", Document: "12345678900"}},
			Paging: domain.Paging{Total: 1},
		}, nil)

		_, err := service.Create(context.Background(), partner)

		require.Error(t, err)
	})

	t.Run("checks duplicates against active and inactive partners alike", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mock_domain.NewMockPartnerRepository(ctrl)
		service := NewPartnerService(mockRepo)

		partner := domain.Partner{Document: "12345678900", DocumentType: domain.CPF}

		mockRepo.EXPECT().Search(gomock.Any(), domain.PartnerFilters{
			DocumentExact: []string{"12345678900"},
			PagingFilter:  domain.PagingFilter{Limit: 1, Offset: 0},
		}).DoAndReturn(func(_ context.Context, filters domain.PartnerFilters) (domain.PagingResult[domain.Partner], error) {
			assert.Nil(t, filters.Active)
			return domain.PagingResult[domain.Partner]{
				Result: []domain.Partner{{PartnerID: "partner-1", Active: false}},
				Paging: domain.Paging{Total: 1},
			}, nil
		})

		_, err := service.Create(context.Background(), partner)

		require.Error(t, err)
	})

	t.Run("creates the partner when the document is not taken", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mock_domain.NewMockPartnerRepository(ctrl)
		service := NewPartnerService(mockRepo)

		partner := domain.Partner{PartnerID: "partner-1", Document: "12345678900", DocumentType: domain.CPF}

		mockRepo.EXPECT().Search(gomock.Any(), domain.PartnerFilters{
			DocumentExact: []string{"12345678900"},
			PagingFilter:  domain.PagingFilter{Limit: 1, Offset: 0},
		}).Return(domain.PagingResult[domain.Partner]{}, nil)
		mockRepo.EXPECT().Create(gomock.Any(), partner).Return("partner-1", nil)

		partnerID, err := service.Create(context.Background(), partner)

		require.NoError(t, err)
		assert.Equal(t, "partner-1", partnerID)
	})
}
