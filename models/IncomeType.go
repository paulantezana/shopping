package models

import "time"

type IncomeType struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	Description string `json:"description"`
	CompanyId   uint   `json:"company_id"`

	CreatedAt     time.Time `json:"-"`
	UpdatedAt     time.Time `json:"-"`
	CreatedUserId uint      `json:"-"`
	UpdatedUserId uint      `json:"-"`
	State         bool      `json:"state" gorm:"default: true"`
}
