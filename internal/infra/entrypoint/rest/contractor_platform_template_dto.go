package rest

import "github.com/icrxz/crm-api-core/internal/domain"

type ContractorPlatformTemplateDTO struct {
	URL           string `json:"url"`
	LoginName     string `json:"username"`
	LoginPassword string `json:"password"`
}

func mapContractorPlatformTemplateDTOToContractorPlatformTemplate(contractorPlatformTemplateDTO ContractorPlatformTemplateDTO) (domain.ContractorPlatformTemplate, error) {
	return domain.NewContractorPlatformTemplate(
		contractorPlatformTemplateDTO.URL,
		contractorPlatformTemplateDTO.LoginName,
		contractorPlatformTemplateDTO.LoginPassword,
	)
}
