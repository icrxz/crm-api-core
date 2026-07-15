package database

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/icrxz/crm-api-core/internal/domain"
)

// JSONMap adapts a map[string]any to a Postgres jsonb column.
type JSONMap map[string]any

func (m JSONMap) Value() (driver.Value, error) {
	if m == nil {
		return "{}", nil
	}

	data, err := json.Marshal(map[string]any(m))
	if err != nil {
		return nil, err
	}

	return string(data), nil
}

func (m *JSONMap) Scan(src any) error {
	if src == nil {
		*m = JSONMap{}
		return nil
	}

	var data []byte
	switch value := src.(type) {
	case []byte:
		data = value
	case string:
		data = []byte(value)
	default:
		return fmt.Errorf("unsupported type for JSONMap: %T", src)
	}

	return json.Unmarshal(data, m)
}

type CaseHistoryDTO struct {
	HistoryID string    `db:"history_id"`
	CaseID    string    `db:"case_id"`
	EventName string    `db:"event_name"`
	AuthorID  string    `db:"author_id"`
	OldValues JSONMap   `db:"old_values"`
	NewValues JSONMap   `db:"new_values"`
	CreatedAt time.Time `db:"created_at"`
}

func mapCaseHistoryToCaseHistoryDTO(history domain.CaseHistory) CaseHistoryDTO {
	return CaseHistoryDTO{
		HistoryID: history.HistoryID,
		CaseID:    history.CaseID,
		EventName: history.EventName,
		AuthorID:  history.AuthorID,
		OldValues: history.OldValues,
		NewValues: history.NewValues,
		CreatedAt: history.CreatedAt,
	}
}

func mapCaseHistoryDTOToCaseHistory(historyDTO CaseHistoryDTO) domain.CaseHistory {
	return domain.CaseHistory{
		HistoryID: historyDTO.HistoryID,
		CaseID:    historyDTO.CaseID,
		EventName: historyDTO.EventName,
		AuthorID:  historyDTO.AuthorID,
		OldValues: historyDTO.OldValues,
		NewValues: historyDTO.NewValues,
		CreatedAt: historyDTO.CreatedAt,
	}
}

func mapCaseHistoryDTOsToCaseHistories(historyDTOs []CaseHistoryDTO) []domain.CaseHistory {
	histories := make([]domain.CaseHistory, 0, len(historyDTOs))
	for _, historyDTO := range historyDTOs {
		histories = append(histories, mapCaseHistoryDTOToCaseHistory(historyDTO))
	}

	return histories
}
