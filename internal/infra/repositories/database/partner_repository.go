package database

import (
	"context"

	"github.com/icrxz/crm-api-core/internal/domain"
	"github.com/jmoiron/sqlx"
)

type providerDatabase struct {
	client *sqlx.DB
}

func NewProviderRepository(client *sqlx.DB) domain.ProviderRepository {
	return &providerDatabase{
		client: client,
	}
}

func (p *providerDatabase) Create(ctx context.Context, provider domain.Provider) (string, error) {
	panic("unimplemented")
}

func (p *providerDatabase) Delete(ctx context.Context, providerID string) error {
	panic("unimplemented")
}

func (p *providerDatabase) GetByID(ctx context.Context, providerID string) (*domain.Provider, error) {
	panic("unimplemented")
}

func (p *providerDatabase) Search(ctx context.Context, filters map[string]string) ([]domain.Provider, error) {
	panic("unimplemented")
}

func (p *providerDatabase) Update(ctx context.Context, providerToUpdate domain.User) error {
	panic("unimplemented")
}
