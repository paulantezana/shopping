package models

type LabelLocation struct {
	ID           uint   `json:"id" gorm:"primary_key"`
	Denomination string `json:"denomination"`
	LevelNumber  uint   `json:"level_number"`
	State        bool   `json:"state" gorm:"default:'true'"`

	CountryID uint `json:"country_id"`
}
