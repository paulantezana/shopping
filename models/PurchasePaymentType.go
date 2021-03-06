package models

type PurchasePaymentType struct {
	ID            uint    `json:"id" gorm:"primaryKey"`
	PurchaseId    uint    `json:"purchase_id"`
	PaymentTypeId uint    `json:"payment_type_id"`
	Total         float64 `json:"total" gorm:"default: 0"`
}
