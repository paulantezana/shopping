package models

import "time"

// UserRoleAuthorization
type UserRoleAuthorization struct {
	ID                 uint `json:"id" gorm:"primaryKey"`
	AppAuthorizationId uint `json:"app_authorization_id"`
	UserRoleId         uint `json:"user_role_id"`

	CreatedAt     time.Time `json:"-"`
	UpdatedAt     time.Time `json:"-"`
	CreatedUserId uint      `json:"-"`
	UpdatedUserId uint      `json:"-"`
	State         bool      `json:"state" gorm:"default: true"`
}
