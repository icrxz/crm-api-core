package builder

import (
	"testing"

	"github.com/icrxz/crm-api-core/internal/domain"
	"github.com/stretchr/testify/require"
)

func TestAssurantBuilder_BuildCase_SetsCategoryMetadata(t *testing.T) {
	columnsIndex := map[string]int{
		"Defeito Reclamado": 0,
		"Número Sinistro":   1,
	}

	builder := NewAssurantBuilder(columnsIndex, "author-1", "Assurant", "d+")

	row := []string{"não liga", "SIN-ASSURANT-001"}
	contractors := []domain.Contractor{{ContractorID: "contractor-assurant"}}

	crmCase, err := builder.BuildCase(row, contractors, "customer-1", 1)

	require.NoError(t, err)
	require.Equal(t, "d+", crmCase.Metadata["category"])
}
