package models

type InvoiceCreditDebit struct {
	ID                    uint `json:"id" gorm:"primary_key"`
	InvoiceId             uint `json:"invoice_id" gorm:"unique; not null"`
	InvoiceParentId       uint `json:"invoice_parent_id" gorm:"not null"`
	UtilCreditDebitTypeId uint `json:"util_credit_debit_type_id" gorm:"not null"`
}
