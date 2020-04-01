package models

import "time"

type CompanyLocal struct {
	ID          uint   `json:"id" gorm:"primary_key"`
	Name        string `json:"name"`
	Phone       string `json:"phone"`
	Address     string `json:"address"`
	Description string `json:"description"`

	UtilGeographicalLocationId uint `json:"util_geographical_location_id"`
	CompanyId                  uint `json:"company_id"`

	CreatedAt     time.Time `json:"-"`
	UpdatedAt     time.Time `json:"-"`
	CreatedUserId uint      `json:"-"`
	UpdatedUserId uint      `json:"-"`
	State         bool      `json:"state"`
}
