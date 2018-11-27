package models

import "time"

type BranchOffice struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Primary   bool      `json:"primary"`
	Address   string    `json:"address"`
	State     bool      `json:"state" gorm:"default:'true'"`

	GeographicLocationID uint `json:"geographic_location_id"`
}
