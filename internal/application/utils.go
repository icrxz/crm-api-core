package application

import "fmt"

func ParseDocument(document string) string {
	switch len(document) {
	case 11:
		return parseCPF(document)
	case 14:
		return parseCNPJ(document)
	default:
		return document
	}
}

func parseCPF(cpf string) string {
	return fmt.Sprintf("%s.%s.%s-%s", cpf[:3], cpf[3:6], cpf[6:9], cpf[9:11])
}

func parseCNPJ(cnpj string) string {
	return fmt.Sprintf("%s.%s.%s/%s-%s", cnpj[:2], cnpj[2:5], cnpj[5:8], cnpj[8:12], cnpj[12:14])
}
