package models

import "time"

// UserForgot
type UserForgot struct {
	ID        uint   `json:"id" gorm:"primaryKey"`
	Used      bool   `json:"used" gorm:"default: false"`
	SecretKey string `json:"secret_key" gorm:"type:varchar(128); not null"`

	UserId uint `json:"user_id"`

	CreatedAt     time.Time `json:"-"`
	UpdatedAt     time.Time `json:"-"`
	CreatedUserId uint      `json:"-"`
	UpdatedUserId uint      `json:"-"`
	State         bool      `json:"state" gorm:"default: true"`
}
