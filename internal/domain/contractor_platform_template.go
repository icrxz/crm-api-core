package domain

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type ContractorPlatformTemplate struct {
	ContractorPlatformTemplateID string
	URL                          string
	LoginName                    string
	LoginPassword                string
	Fields                       map[string]map[string]string
}

func NewContractorPlatformTemplate(url, loginName, loginPassword string) (ContractorPlatformTemplate, error) {
	contractorPlatformTemplateID, err := uuid.NewRandom()
	if err != nil {
		return ContractorPlatformTemplate{}, err
	}

	encryptedLogin, err := encryptData(loginName)
	if err != nil {
		return ContractorPlatformTemplate{}, err
	}
	encryptedPassword, err := encryptData(loginPassword)
	if err != nil {
		return ContractorPlatformTemplate{}, err
	}

	return ContractorPlatformTemplate{
		ContractorPlatformTemplateID: contractorPlatformTemplateID.String(),
		URL:                          url,
		LoginName:                    encryptedLogin,
		LoginPassword:                encryptedPassword,
	}, nil
}

func encryptData(input string) (string, error) {
	encryptedInput, err := bcrypt.GenerateFromPassword([]byte(input), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(encryptedInput), nil
}
