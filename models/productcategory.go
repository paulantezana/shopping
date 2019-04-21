package models

import "time"

type ProductCategory struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	State     bool      `json:"state" gorm:"default:'true'"`

	MainRelationship uint `json:"main_relationship"`
	ProductID        uint `json:"product_id"`
	CategoryID       uint `json:"category_id"`
}
