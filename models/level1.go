package models

type Level1 struct {
	ID    uint   `json:"id" gorm:"primary_key"`
	Name  string `json:"name"`
	State bool   `json:"state" gorm:"default:'true'"`

	CountryID uint `json:"country_id"`
}
