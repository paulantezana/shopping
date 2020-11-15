package models

type GroupAttribute struct {
	ID    uint   `json:"id" gorm:"primary_key"`
	Name  string `json:"name" gorm:"type:varchar(12)"`
	State bool   `json:"state" gorm:"default: true"`
}
