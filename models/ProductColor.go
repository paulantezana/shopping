package models

import "time"

type ProductColor struct {
	ID uint `json:"id" gorm:"primary_key"`

	Description      string `json:"description"  gorm:"not null"`
	ShortDescription string `json:"short_description" gorm:"default: ''"`
	Color            string `json:"color"  gorm:"not null"`

	ProductId uint `json:"product_id"  gorm:"not null"`

	CreatedAt     time.Time `json:"-"`
	UpdatedAt     time.Time `json:"-"`
	CreatedUserId uint      `json:"-"`
	UpdatedUserId uint      `json:"-"`
	State         bool      `json:"state"`
}
