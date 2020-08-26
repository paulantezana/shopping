package models

import "time"

type CompanySalePoint struct {
	ID          uint   `json:"id" gorm:"primary_key"`
	Description string `json:"description" gorm:"type:varchar(128)"`

	CompanyLocalId uint `json:"company_local_id"`
	CompanyId      uint `json:"company_id"`

	CreatedAt     time.Time `json:"-"`
	UpdatedAt     time.Time `json:"-"`
	CreatedUserId uint      `json:"-"`
	UpdatedUserId uint      `json:"-"`
	State         bool      `json:"state" gorm:"default: true"`
}
