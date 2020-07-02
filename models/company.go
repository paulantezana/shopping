package models

import "time"

type Company struct {
	ID               uint   `json:"id" gorm:"primary_key"`
	DocumentNumber   string `json:"document_number"`
	SocialReason     string `json:"social_reason"`
	CommercialReason string `json:"commercial_reason"`
	Representative   string `json:"representative"`
	Logo             string `json:"logo"`
	Email            string `json:"email"`
	Phone            string `json:"phone"`
	Address          string `json:"address"`
	Color1           string `json:"color_1"`
	Color2           string `json:"color_2"`

	UtilGeographicalLocationId uint `json:"util_geographical_location_id"`

	CreatedAt     time.Time `json:"-"`
	UpdatedAt     time.Time `json:"-"`
	CreatedUserId uint      `json:"-"`
	UpdatedUserId uint      `json:"-"`
	State         bool      `json:"state" gorm:"default: true"`
}
