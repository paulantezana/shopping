package models

import "time"

// User user data
type User struct {
	ID uint `json:"id" gorm:"primary_key"`

	CompanyId uint `json:"company_id"`

	CreatedAt     time.Time `json:"-"`
	UpdatedAt     time.Time `json:"-"`
	CreatedUserId uint      `json:"-"`
	UpdatedUserId uint      `json:"-"`
	State         bool      `json:"state"`
}
