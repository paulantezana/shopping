package models

import "time"

type ProductPresentation struct {
	ID uint `json:"id" gorm:"primary_key"`

	Barcode     string `json:"barcode"`
	Description string `json:"description" gorm:"default: ''"`

	SaleUnitPrice1    float64 `json:"sale_unit_price_1"`
	SaleUnitPrice2    float32 `json:"sale_unit_price_2"`
	SaleUnitValue1    float32 `json:"sale_unit_value_1"`
	SaleUnitValue2    float32 `json:"sale_unit_value_2"`
	PurchaseUnitPrice float32 `json:"purchase_price"`
	PurchaseUnitValue float32 `json:"purchase_unit_value"`

	ProductColorId  uint `json:"product_color_id"`
	ProductSizeId   uint `json:"product_size_id"`
	ProductMaterial uint `json:"product_material"`

	ProductId uint `json:"product_id"  gorm:"not null"`

	CreatedAt     time.Time `json:"-"`
	UpdatedAt     time.Time `json:"-"`
	CreatedUserId uint      `json:"-"`
	UpdatedUserId uint      `json:"-"`
	State         bool      `json:"state"`
}
