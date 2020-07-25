package models

import "time"

type CompanyWareHouse struct {
	ID             uint `json:"id" gorm:"primary_key"`
	CompanyLocalId uint `json:"company_local_id"`

	CreatedAt     time.Time `json:"-"`
	UpdatedAt     time.Time `json:"-"`
	CreatedUserId uint      `json:"-"`
	UpdatedUserId uint      `json:"-"`
	State         bool      `json:"state" gorm:"default: true"`
}