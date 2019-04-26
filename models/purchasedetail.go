package models

type PurchaseDetail struct {
	ID           uint    `json:"id" gorm:"primary_key"`
	Description  string  `json:"description"`
	Quantity     float32 `json:"quantity"`
	UnitQuantity float32 `json:"unit_quantity"`
	UnitPrice    float32 `json:"unit_price"`
	Discount     float32 `json:"discount"`
	Total        float32 `json:"total"`
	State        bool    `json:"state" gorm:"default:'true'"`

	PresentationID uint `json:"presentation_id"`
	ProductID      uint `json:"product_id"`
	PurchaseID     uint `json:"purchase_id"`
}
