package models

// UtilCreditDebitType --
type UtilCreditDebitType struct {
	ID                 uint   `json:"id" gorm:"primary_key"`
	Code               string `json:"code" gorm:"type:varchar(12)"`
	Description        string `json:"description" gorm:"type:varchar(128)"`
	UtilDocumentTypeId uint   `json:"util_document_type_id"`
	State              bool   `json:"state" gorm:"default: true"`
}
