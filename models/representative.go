package models

type Representative struct {
	ID             uint   `json:"id" gorm:"primary_key"`
	Document       string `json:"document"`
	NumberDocument string `json:"number_document"`
	Name           string `json:"name"`
	Position       string `json:"position"`
	State          bool   `json:"state" gorm:"default:'true'"`

	GeneralDataID uint `json:"general_data_id"`
}
