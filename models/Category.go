package models

import "time"

// Category --
type Category struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	Name        string `json:"name" gorm:"type:varchar(64)"`
	Description string `json:"description" gorm:"type:varchar(128)"`

	CompanyId uint `json:"company_id"`
	ParentId  uint `json:"parent_id"`

	CreatedAt     time.Time `json:"-"`
	UpdatedAt     time.Time `json:"-"`
	CreatedUserId uint      `json:"-"`
	UpdatedUserId uint      `json:"-"`
	State         bool      `json:"state" gorm:"default: true"`
}
