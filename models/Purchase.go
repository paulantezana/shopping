package models

import "time"

type Purchase struct {
	ID                uint      `json:"id" gorm:"primary_key"`
	DateOfIssue       time.Time `json:"date_of_issue" gorm:"not null"`
	DateOfPurchase    time.Time `json:"date_of_purchase"`
	Invoice           string    `json:"invoice"` // Folio
	Guide             string    `json:"guide" `
	Observation       string    `json:"observation" gorm:"default: ''"`
	CancelObservation string    `json:"cancel_observation" gorm:"default: ''"`
	TotalUnaffected   float64   `json:"total_unaffected" gorm:"default: 0.00"`
	TotalTaxed        float64   `json:"total_taxed" gorm:"default: 0.00"`
	TotalIgv          float64   `json:"total_igv" gorm:"default: 0.00"`
	Total             float64   `json:"total" gorm:"default: 0.00"`
	TotalInLetter     string    `json:"total_in_letter"`

	ProviderId         uint `json:"provider_id"`
	CompanyWareHouseId uint `json:"company_ware_house_id"`
	UtilCurrencyTypeId uint `json:"util_currency_type_id"`
	UtilDocumentTypeId uint `json:"util_document_type_id"`
	CompanyId          uint `json:"company_id"`
	UserId             uint `json:"user_id"`

	CreatedAt     time.Time `json:"-"`
	UpdatedAt     time.Time `json:"-"`
	CreatedUserId uint      `json:"-"`
	UpdatedUserId uint      `json:"-"`
	State         bool      `json:"state" gorm:"default: true"`
}
