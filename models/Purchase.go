package models

import "time"

type Purchase struct {
	ID              uint      `json:"id" gorm:"primary_key"`
	DateOfIssue     time.Time `json:"date_of_issue" gorm:"not null"`
	DateOfPurchase  time.Time `json:"date_of_purchase"`
	Serie           string    `json:"serie"`
	Number          string    `json:"number"`
	PdfFormat       string    `json:"pdf_format" gorm:"default: ''"`
	Guide           string    `json:"guide" `
	Observation     string    `json:"observation" gorm:"default: ''"`
	TotalUnaffected float64   `json:"total_unaffected" gorm:"default: 0.00"`
	TotalTaxed      float64   `json:"total_taxed" gorm:"default: 0.00"`
	TotalIgv        float64   `json:"total_igv" gorm:"default: 0.00"`
	Total           float64   `json:"total" gorm:"default: 0.00"`

	ProviderId         uint `json:"provider_id"`
	CompanyWareHouseId uint `json:"company_ware_house_id"`
	UtilCurrencyTypeId uint `json:"util_currency_type_id"`
	UtilDocumentTypeId   uint `json:"util_document_type_id"`

	CreatedAt     time.Time `json:"-"`
	UpdatedAt     time.Time `json:"-"`
	CreatedUserId uint      `json:"-"`
	UpdatedUserId uint      `json:"-"`
	State         bool      `json:"state" gorm:"default: true"`
}
