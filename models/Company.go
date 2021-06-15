package models

import "time"

// Company model
type Company struct {
	ID               uint   `json:"id" gorm:"primaryKey"`
	DocumentNumber   string `json:"document_number" gorm:"unique; not null; type:varchar(32)"`
	SocialReason     string `json:"social_reason" gorm:"not null; type:varchar(255)"`
	CommercialReason string `json:"commercial_reason" gorm:"type:varchar(255)"`
	Representative   string `json:"representative" gorm:"type:varchar(128)"`
	Email            string `json:"email" gorm:"type:varchar(64)"`
	Logo             string `json:"logo" gorm:"type:varchar(128)"`
	LogoLarge        string `json:"logo_large" gorm:"type:varchar(128)"`
	Phone            string `json:"phone" gorm:"type:varchar(32)"`
	Address          string `json:"address" gorm:"type:varchar(255)"`
	Color1           string `json:"color_1" gorm:"type:varchar(32)"`
	Color2           string `json:"color_2" gorm:"type:varchar(32)"`

	UtilGeographicalLocationId uint `json:"util_geographical_location_id"`

	CreatedAt     time.Time `json:"-"`
	UpdatedAt     time.Time `json:"-"`
	CreatedUserId uint      `json:"-"`
	UpdatedUserId uint      `json:"-"`
	State         bool      `json:"state" gorm:"default: true"`
}
