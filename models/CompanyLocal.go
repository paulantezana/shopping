package models

import (
	"time"
)

type CompanyLocal struct {
	ID               uint   `json:"id" gorm:"primary_key"`
	SocialReason     string `json:"social_reason"`
	CommercialReason string `json:"commercial_reason"`
	Phone            string `json:"phone"`
	Address          string `json:"address"`
	Description      string `json:"description"`

	UtilGeographicalLocationId uint `json:"util_geographical_location_id"`
	CompanyId                  uint `json:"company_id"`

	CreatedAt     time.Time `json:"-"`
	UpdatedAt     time.Time `json:"-"`
	CreatedUserId uint      `json:"-"`
	UpdatedUserId uint      `json:"-"`
	State         bool      `json:"state" gorm:"default: true"`

	CompanySeries []CompanySerie `json:"company_series"`
}
