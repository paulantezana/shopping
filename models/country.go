package models

type Country struct {
	ID    uint   `json:"id" gorm:"primary_key"`
	Name  string `json:"name"`
	Code  string `json:"code"`
	State bool   `json:"state" gorm:"default:'true'"`
}
