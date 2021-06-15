package models

import (
	"time"
)

type CompanyLocal struct {
	ID               uint   `json:"id" gorm:"primaryKey"`
	SocialReason     string `json:"social_reason" gorm:"type:varchar(255)"`
	CommercialReason string `json:"commercial_reason" gorm:"type:varchar(255)"`
	Phone            string `json:"phone" gorm:"type:varchar(128)"`
	Address          string `json:"address" gorm:"type:varchar(255)"`
	Description      string `json:"description" gorm:"type:varchar(128)"`

	UtilGeographicalLocationId uint `json:"util_geographical_location_id"`
	CompanyId                  uint `json:"company_id"`

	CreatedAt     time.Time `json:"-"`
	UpdatedAt     time.Time `json:"-"`
	CreatedUserId uint      `json:"-"`
	UpdatedUserId uint      `json:"-"`
	State         bool      `json:"state" gorm:"default: true"`
}
