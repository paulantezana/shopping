package models

import "time"

type Attribute struct {
	ID        uint   `json:"id" gorm:"primaryKey"`
	Name      string `json:"name" gorm:"type:varchar(12)"`
	ShortName string `json:"short_name" gorm:"type:varchar(12)"`

	GroupAttributeId int `json:"group_attribute_id"`

	CreatedAt     time.Time `json:"-"`
	UpdatedAt     time.Time `json:"-"`
	CreatedUserId uint      `json:"-"`
	UpdatedUserId uint      `json:"-"`
	State         bool      `json:"state" gorm:"default: true"`
}
