package utils

import (
	"regexp"
	"strings"
)

// IsCPFValid valida um CPF (simples checagem de formato e dígitos verificadores)
func IsCPFValid(cpf string) bool {
	cpf = strings.ReplaceAll(cpf, ".", "")
	cpf = strings.ReplaceAll(cpf, "-", "")

	if len(cpf) != 11 {
		return false
	}

	// TODO: Implementar algoritmo completo de validação de CPF
	// Por enquanto, apenas checa se é numérico
	match, _ := regexp.MatchString("^[0-9]+$", cpf)
	return match
}

// IsCNPJValid valida um CNPJ
func IsCNPJValid(cnpj string) bool {
	cnpj = strings.ReplaceAll(cnpj, ".", "")
	cnpj = strings.ReplaceAll(cnpj, "/", "")
	cnpj = strings.ReplaceAll(cnpj, "-", "")

	if len(cnpj) != 14 {
		return false
	}

	match, _ := regexp.MatchString("^[0-9]+$", cnpj)
	return match
}
