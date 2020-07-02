package models

import "time"

// User user data
type User struct {
	ID       uint   `json:"id" gorm:"primary_key"`
	UserName string `json:"user_name" gorm:"type:varchar(64); not null"`
	Password string `json:"password" gorm:"type:varchar(64); not null"`
	Avatar   string `json:"avatar" gorm:"type:varchar(255)"`
	Email    string `json:"email" gorm:"type:varchar(64); not null"`
	Freeze   bool   `json:"-"`

	CompanyId  uint `json:"company_id"`
	UserRoleId uint `json:"user_role_id"`

	CreatedAt     time.Time `json:"-"`
	UpdatedAt     time.Time `json:"-"`
	CreatedUserId uint      `json:"-"`
	UpdatedUserId uint      `json:"-"`
	State         bool      `json:"state" gorm:"default: true"`
}
