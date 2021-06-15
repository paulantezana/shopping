package models

import "time"

type CompanySalePoint struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	Description string `json:"description" gorm:"type:varchar(128)"`

	CompanyLocalId uint `json:"company_local_id"`
	CompanyId      uint `json:"company_id"`

	CreatedAt     time.Time `json:"-"`
	UpdatedAt     time.Time `json:"-"`
	CreatedUserId uint      `json:"-"`
	UpdatedUserId uint      `json:"-"`
	State         bool      `json:"state" gorm:"default: true"`
	
	CompanySalePointSeries []CompanySalePointSerie `json:"company_sale_point_series"`
}
