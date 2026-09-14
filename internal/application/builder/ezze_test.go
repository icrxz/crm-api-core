package builder

import (
	"testing"

	"github.com/icrxz/crm-api-core/internal/domain"
	"github.com/stretchr/testify/require"
)

func TestEzzeBuilder_BuildCase_SetsCategoryMetadata(t *testing.T) {
	columnsIndex := map[string]int{
		"Ticket": 0,
	}

	builder := NewEzzeBuilder(columnsIndex, "author-1", "Ezze Seguros", "furniture")

	row := []string{"TICKET-001"}
	contractors := []domain.Contractor{{ContractorID: "contractor-ezze"}}

	crmCase, err := builder.BuildCase(row, contractors, "customer-1", 1)

	require.NoError(t, err)
	require.Equal(t, "furniture", crmCase.Metadata["category"])
}
