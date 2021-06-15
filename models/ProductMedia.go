package models

import "time"

type ProductMedia struct {
	ID uint `json:"id" gorm:"primaryKey"`

	ProductId      uint `json:"product_id" gorm:"not null"`
	ProductColorId uint `json:"product_color_id" gorm:"not null"`

	CreatedAt     time.Time `json:"-"`
	UpdatedAt     time.Time `json:"-"`
	CreatedUserId uint      `json:"-"`
	UpdatedUserId uint      `json:"-"`
	State         bool      `json:"state"`
}
