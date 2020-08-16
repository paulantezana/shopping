package models

import "time"

// UserRole
type UserRole struct {
	ID          uint   `json:"id" gorm:"primary_key"`
	Description string `json:"description" gorm:"type:varchar(64)"`
	CompanyId   uint   `json:"company_id"`

	CreatedAt     time.Time `json:"-"`
	UpdatedAt     time.Time `json:"-"`
	CreatedUserId uint      `json:"-"`
	UpdatedUserId uint      `json:"-"`
	State         bool      `json:"state" gorm:"default: true"`

	UserRoleAuthorizations []UserRoleAuthorization `json:"user_role_authorizations"`
}
