package builder

import (
	"testing"

	"github.com/icrxz/crm-api-core/internal/domain"
	"github.com/stretchr/testify/require"
)

func TestDefaultBuilder_BuildCase_SetsCategoryMetadata(t *testing.T) {
	columnsIndex := map[string]int{
		"Sinistro":  0,
		"Descrição": 1,
		"Documento": 2,
		"Nome":      3,
		"Sobrenome": 4,
		"Cidade":    5,
		"Estado":    6,
		"Marca":     7,
		"Modelo":    8,
		"Valor":     9,
	}

	builder := NewDefaultBuilder(columnsIndex, "author-1", "Seguradora ABC", "d+")

	row := []string{"SIN-DEFAULT-001", "quebrado", "111.111.111-11", "Maria", "Silva", "Curitiba", "PR", "Brastemp", "BRM45", "100"}
	contractors := []domain.Contractor{{ContractorID: "contractor-default"}}

	crmCase, err := builder.BuildCase(row, contractors, "customer-1", 1)

	require.NoError(t, err)
	require.Equal(t, "d+", crmCase.Metadata["category"])
}
