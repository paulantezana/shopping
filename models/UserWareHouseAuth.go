package models

import "time"

type UserWareHouseAuth struct {
	ID                 uint `json:"id" gorm:"primaryKey"`
	CompanyWareHouseId uint `json:"company_ware_house_id"`
	UserId             uint `json:"user_id"`

	CreatedAt     time.Time `json:"-"`
	UpdatedAt     time.Time `json:"-"`
	CreatedUserId uint      `json:"-"`
	UpdatedUserId uint      `json:"-"`
	State         bool      `json:"state" gorm:"default: true"`
}
