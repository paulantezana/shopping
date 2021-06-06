package models

type InvoiceCustomer struct {
	ID             uint   `json:"id" gorm:"primary_key"`
	DocumentNumber string `json:"document_number" gorm:"default: ''"`
	SocialReason   string `json:"social_reason" gorm:"default: ''"`
	FiscalAddress  string `json:"fiscal_address" gorm:"default: ''"`
	Email          string `json:"email" gorm:"default: ''"`
	Phone          string `json:"phone" gorm:"default: ''"`
	EmailSend      bool   `json:"email_send" gorm:"default: false"`
	PhoneSend      bool   `json:"phone_send" gorm:"default: false"`

	InvoiceId                  uint `json:"invoice_id" gorm:"unique; not null"`
	CustomerId uint `json:"customer_id"`
	UtilIdentityDocumentTypeId uint `json:"util_identity_document_type_id" gorm:"not null"`
}
