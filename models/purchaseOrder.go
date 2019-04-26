package models

import "time"

type PurchaseOrder struct {
	ID           uint      `json:"id" gorm:"primary_key"`
	Series       string    `json:"series"`
	Correlative  string    `json:"correlative"`
	Date         time.Time `json:"date"`
	DeliveryTerm time.Time `json:"delivery_term"`
	Observation  string    `json:"observation"`
	Address      string    `json:"address"`
	State        bool      `json:"state" gorm:"default:'true'"`

	GeographicLocationID uint `json:"geographic_location_id"`
	DocumentTypeID       uint `json:"document_type_id"`
	PurchaseID           uint `json:"purchase_id"`
}
