package database

import (
	"context"

	"github.com/icrxz/crm-api-core/internal/domain"
	"github.com/jmoiron/sqlx"
)

type caseHistoryRepository struct {
	client *sqlx.DB
}

func NewCaseHistoryRepository(client *sqlx.DB) domain.CaseHistoryRepository {
	return &caseHistoryRepository{
		client: client,
	}
}

func (r *caseHistoryRepository) Create(ctx context.Context, history domain.CaseHistory) error {
	historyDTO := mapCaseHistoryToCaseHistoryDTO(history)

	_, err := executor(ctx, r.client).NamedExecContext(
		ctx,
		"INSERT INTO case_history "+
			"(history_id, case_id, event_name, author_id, old_values, new_values, created_at) "+
			"VALUES "+
			"(:history_id, :case_id, :event_name, :author_id, :old_values, :new_values, :created_at)",
		historyDTO,
	)

	return err
}

func (r *caseHistoryRepository) GetByCaseID(ctx context.Context, caseID string) ([]domain.CaseHistory, error) {
	if caseID == "" {
		return nil, domain.NewValidationError("caseID is required", nil)
	}

	var historyDTOs []CaseHistoryDTO
	err := executor(ctx, r.client).SelectContext(
		ctx,
		&historyDTOs,
		"SELECT * FROM case_history WHERE case_id=$1 ORDER BY created_at ASC",
		caseID,
	)
	if err != nil {
		return nil, err
	}

	return mapCaseHistoryDTOsToCaseHistories(historyDTOs), nil
}
