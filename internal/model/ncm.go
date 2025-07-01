package model

import "fmt"

type NCM struct {
	ID            string `json:"id_ncm"`
	Code          string `json:"codigo"`
	CodeNoSimbols string
	Description   string `json:"descricao"`
	InitialDate   string `json:"data_inicio"`
	FinalDate     string `json:"data_fim"`
	TypeYearIni   string `json:"tipo_ato_ini"`
	NumberAtoIni  string `json:"numero_ato_ini"`
	YearAtoIni    string `json:"ano_ato_ini"`
}

func NewNCM(id, code, description, initialDate, finalDate, typeYearIni, numberAtoIni, yearAtoIni string) *NCM {
	return &NCM{
		ID:            id,
		Code:          code,
		CodeNoSimbols: code,
		Description:   description,
		InitialDate:   initialDate,
		FinalDate:     finalDate,
		TypeYearIni:   typeYearIni,
		NumberAtoIni:  numberAtoIni,
		YearAtoIni:    yearAtoIni,
	}
}

func ValidateNCM(ncm *NCM) error {
	if ncm == nil {
		return fmt.Errorf("ncm cannot be nil")
	}

	if ncm.ID == "" {
		return fmt.Errorf("id cannot be null or empty")
	}

	if ncm.Code == "" {
		return fmt.Errorf("code cannot be null or empty")
	}

	if ncm.Description == "" {
		return fmt.Errorf("description cannot be null or empty")
	}

	return nil
}
