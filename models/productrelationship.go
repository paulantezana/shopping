package models

import "time"

type ProductRelationship struct {
	ID               uint      `json:"id" gorm:"primary_key"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	TypeRelationship string    `json:"type_relationship"`
	Position         uint      `json:"position"`
	State            bool      `json:"state" gorm:"default:'true'"`

	ProductRelationshipID uint `json:"product_relationship_id"`
	ProductID             uint `json:"product_id"`
}
