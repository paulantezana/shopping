package models

type PurchaseConf struct {
	ID                     uint `json:"id" gorm:"primary_key"`
	Decimals               uint `json:"decimals"`
	ValidateDocumentNumber bool `json:"validate_document_number"`
	AutoSendEmail          bool `json:"auto_send_email"`
	Email                  bool `json:"email"`
}
