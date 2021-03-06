package models

import "time"

type CompanyWareHouse struct {
	ID             uint   `json:"id" gorm:"primaryKey"`
	CompanyLocalId uint   `json:"company_local_id"`
	CompanyId      uint   `json:"company_id"`
	Description    string `json:"description" gorm:"type:varchar(128)"`

	CreatedAt     time.Time `json:"-"`
	UpdatedAt     time.Time `json:"-"`
	CreatedUserId uint      `json:"-"`
	UpdatedUserId uint      `json:"-"`
	State         bool      `json:"state" gorm:"default: true"`
}
