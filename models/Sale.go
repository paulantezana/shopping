package models

import "time"

type Sale struct {
	ID uint `json:"id" gorm:"primaryKey"`

	DateOfIssue       time.Time `json:"date_of_issue" gorm:"not null"`
	TimeOfIssue       time.Time `json:"time_of_issue" gorm:"not null"`
	DateOfDue         time.Time `json:"date_of_due" gorm:"not null"`
	Serie             string    `json:"serie" gorm:"not null; unique_index:idx_sale_key"`
	Number            uint      `json:"number" gorm:"not null; unique_index:idx_sale_key"`
	ChangeType        string    `json:"change_type" gorm:"default: ''"`
	VehiclePlate      string    `json:"vehicle_plate" gorm:"default: ''"`
	Term              string    `json:"term" gorm:"default: ''"`
	PdfFormat         string    `json:"pdf_format" gorm:"default: ''"`
	Guide             string    `json:"guide" `
	Observation       string    `json:"observation" gorm:"default: ''"`
	CancelObservation string    `json:"cancel_observation" gorm:"default: ''"`

	// Customer
	DocumentNumber             string `json:"document_number" gorm:"default: ''"`
	SocialReason               string `json:"social_reason" gorm:"default: ''"`
	FiscalAddress              string `json:"fiscal_address" gorm:"default: ''"`
	Email                      string `json:"email" gorm:"default: ''"`
	Phone                      string `json:"phone" gorm:"default: ''"`
	EmailSend                  bool   `json:"email_send" gorm:"default: false"`
	PhoneSend                  bool   `json:"phone_send" gorm:"default: false"`
	CustomerId                 uint   `json:"customer_id"`
	UtilIdentityDocumentTypeId uint   `json:"util_identity_document_type_id" gorm:"not null"`
	// End Customer

	TotalFree          float64 `json:"total_free" gorm:"default: 0"`
	TotalExportation   float64 `json:"total_exportation" gorm:"default: 0"`
	TotalDiscount      float64 `json:"total_discount" gorm:"default: 0"`
	TotalExonerated    float64 `json:"total_exonerated" gorm:"default: 0"`
	TotalUnaffected    float64 `json:"total_unaffected" gorm:"default: 0"`
	TotalTaxed         float64 `json:"total_taxed" gorm:"default: 0"`
	TotalIgv           float64 `json:"total_igv" gorm:"default: 0"`
	TotalCharge        float64 `json:"total_charge" gorm:"default: 0"`
	TotalValue         float64 `json:"total_value" gorm:"default: 0"`
	TotalPlasticBagTax float64 `json:"total_plastic_bag_tax" gorm:"default: 0"`
	Total              float64 `json:"total" gorm:"default: 0"`

	UtilDocumentTypeId  uint `json:"util_document_type_id" gorm:"not null; unique_index:idx_sale_key"`
	UtilCurrencyTypeId  uint `json:"util_currency_type_id" gorm:"not null"`
	UtilOperationTypeId uint `json:"util_operation_type_id" gorm:"not null"`
	CompanyId           uint `json:"company_id" gorm:"not null"`
	CompanyWareHouseId  uint `json:"company_ware_house_id"`
	CompanySalePointId  uint `json:"company_sale_point_id"`
	UserId              uint `json:"user_id" gorm:"not null"`

	CreatedAt     time.Time `json:"-"`
	UpdatedAt     time.Time `json:"-"`
	CreatedUserId uint      `json:"-"`
	UpdatedUserId uint      `json:"-"`
	State         bool      `json:"state" gorm:"default: true"`
}
