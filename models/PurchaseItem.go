package models

import "time"

type PurchaseItem struct {
	ID          uint    `json:"id" gorm:"primaryKey"`
	UnitMeasure string  `json:"unit_measure" gorm:"not null"`
	ProductCode string  `json:"product_code" gorm:"not null"`
	Description string  `json:"description" gorm:"not null"`
	Quantity    float64 `json:"quantity" gorm:"not null"`
	UnitPrice   float64 `json:"unit_price" gorm:"not null"`
	Total       float64 `json:"total" gorm:"default: 0"`

	ProductId  uint `json:"product_id" gorm:"not null"`
	PurchaseId uint `json:"purchase_id" gorm:"not null"`

	CreatedAt     time.Time `json:"-"`
	UpdatedAt     time.Time `json:"-"`
	CreatedUserId uint      `json:"-"`
	UpdatedUserId uint      `json:"-"`
	State         bool      `json:"state" gorm:"default: true"`
}
