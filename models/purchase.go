package models

import "time"

type Purchase struct {
	ID               uint      `json:"id" gorm:"primary_key"`
	DocumentNumber   string    `json:"document_number"`
	ProviderName     string    `json:"provider_name"`
	RucDni           string    `json:"ruc_dni"`
	Address          string    `json:"address"`
	PaymentForm      string    `json:"payment_form"`
	Currency         string    `json:"currency"`
	Date             time.Time `json:"date"`
	BillingDate      time.Time `json:"billing_date"`
	ModificationDate time.Time `json:"modification_date"`
	PayDate          time.Time `json:"pay_date"`
	Discount         float32   `json:"discount"`
	TypePurchase     string    `json:"type_purchase"`
	TypeChange       float32   `json:"type_change"`
	Subtotal         float32   `json:"subtotal"`
	Total            float32   `json:"total"`
	Observation      string    `json:"observation"`
	State            bool      `json:"state" gorm:"default:'true'"`

	ProviderID     uint `json:"provider_id"`
	PaymentID      uint `json:"payment_id"`
	PersonalID     uint `json:"personal_id"`
	TypeDocumentID uint `json:"type_document_id"`
	BranchOfficeID uint `json:"branch_office_id"`
}
