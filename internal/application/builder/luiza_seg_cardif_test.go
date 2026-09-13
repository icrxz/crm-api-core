package builder

import (
	"testing"

	"github.com/icrxz/crm-api-core/internal/domain"
	"github.com/stretchr/testify/require"
)

func TestLuizaSegBuilder_BuildCase_SetsCategoryMetadata(t *testing.T) {
	columnsIndex := map[string]int{
		"BASE":     0,
		"SINISTRO": 2,
	}

	const contractorCompanyName = "LuizaSeg"
	builder := NewLuizaSegBuilder(columnsIndex, "author-1", "", "furniture")

	row := []string{"base-col", "Retail", "sinistro-col", "SIN-LUIZASEG-001"}
	contractors := []domain.Contractor{{ContractorID: "contractor-luizaseg", CompanyName: contractorCompanyName}}

	crmCase, err := builder.BuildCase(row, contractors, "customer-1", 1)

	require.NoError(t, err)
	require.Equal(t, "furniture", crmCase.Metadata["category"])
}
