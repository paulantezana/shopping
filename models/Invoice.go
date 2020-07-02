package models

import "time"

type Invoice struct {
	ID uint `json:"id" gorm:"primary_key"`

	DateOfIssue  time.Time `json:"date_of_issue" gorm:"not null"`
	TimeOfIssue  time.Time `json:"time_of_issue" gorm:"not null"`
	DateOfDue    time.Time `json:"date_of_due" gorm:"not null"`
	Serie        string    `json:"serie" gorm:"not null; unique_index:idx_invoice_key"`
	Number       uint      `json:"number" gorm:"not null; unique_index:idx_invoice_key"`
	ChangeType   string    `json:"change_type" gorm:"default: ''"`
	VehiclePlate string    `json:"vehicle_plate" gorm:"default: ''"`
	Term         string    `json:"term" gorm:"default: ''"`
	PdfFormat    string    `json:"pdf_format" gorm:"default: ''"`
	Guide        string    `json:"guide" `
	Observation  string    `json:"observation" gorm:"default: ''"`

	TotalFree          float64 `json:"total_free" gorm:"default: 0.00"`
	TotalExportation   float64 `json:"total_exportation" gorm:"default: 0.00"`
	TotalDiscount      float64 `json:"total_discount" gorm:"default: 0.00"`
	TotalExonerated    float64 `json:"total_exonerated" gorm:"default: 0.00"`
	TotalUnaffected    float64 `json:"total_unaffected" gorm:"default: 0.00"`
	TotalTaxed         float64 `json:"total_taxed" gorm:"default: 0.00"`
	TotalIgv           float64 `json:"total_igv" gorm:"default: 0.00"`
	TotalCharge        float64 `json:"total_charge" gorm:"default: 0.00"`
	TotalValue         float64 `json:"total_value" gorm:"default: 0.00"`
	TotalPlasticBagTax float64 `json:"total_plastic_bag_tax" gorm:"default: 0.00"`
	Total              float64 `json:"total" gorm:"default: 0.00"`

	UtilDocumentTypeId  uint `json:"util_document_type_id" gorm:"not null; unique_index:idx_invoice_key"`
	UtilCurrencyTypeId  uint `json:"util_currency_type_id" gorm:"not null"`
	UtilOperationTypeId uint `json:"util_operation_type_id" gorm:"not null"`
	CompanyId           uint `json:"company_id" gorm:"not null"`
	CompanyLocalId      uint `json:"company_local_id" gorm:"not null; unique_index:idx_invoice_key"`

	CreatedAt     time.Time `json:"-"`
	UpdatedAt     time.Time `json:"-"`
	CreatedUserId uint      `json:"-"`
	UpdatedUserId uint      `json:"-"`
	State         bool      `json:"state" gorm:"default: true"`
}
