package models

import "time"

type UserSalePoint struct {
	ID                 uint `json:"id" gorm:"primaryKey"`
	CompanySalePointId uint `json:"company_sale_point_id"`
	UserId             uint `json:"user_id"`

	CreatedAt     time.Time `json:"-"`
	UpdatedAt     time.Time `json:"-"`
	CreatedUserId uint      `json:"-"`
	UpdatedUserId uint      `json:"-"`
	State         bool      `json:"state" gorm:"default: true"`
}
