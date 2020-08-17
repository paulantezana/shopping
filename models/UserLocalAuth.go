package models

import "time"

type UserLocalAuth struct {
	ID             uint `json:"id" gorm:"primary_key"`
	CompanyLocalId uint `json:"company_local_id"`
	UserId         uint `json:"user_id"`

	CreatedAt     time.Time `json:"-"`
	UpdatedAt     time.Time `json:"-"`
	CreatedUserId uint      `json:"-"`
	UpdatedUserId uint      `json:"-"`
	State         bool      `json:"state" gorm:"default: true"`
}
