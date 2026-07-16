package rest

import (
	"testing"

	"github.com/icrxz/crm-api-core/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestMapCaseFullToCaseFullDTO(t *testing.T) {
	caseFull := domain.CaseFull{
		CaseID: "case-1",
		Queue: domain.Queue{
			QueueID:  "queue-1",
			Name:     "SP Mobile",
			Category: domain.MobileQueueCategory,
			States:   []string{"SP"},
		},
	}

	dto := mapCaseFullToCaseFullDTO(caseFull)

	assert.Equal(t, "case-1", dto.CaseID)
	assert.Equal(t, "queue-1", dto.Queue.QueueID)
	assert.Equal(t, "SP Mobile", dto.Queue.Name)
	assert.Equal(t, []string{"SP"}, dto.Queue.States)
}

func TestMapCaseFullToCaseFullDTO_EmptyQueue(t *testing.T) {
	caseFull := domain.CaseFull{CaseID: "case-1"}

	dto := mapCaseFullToCaseFullDTO(caseFull)

	assert.Empty(t, dto.Queue.QueueID)
}
