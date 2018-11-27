package models

type Level2 struct {
	ID    uint   `json:"id" gorm:"primary_key"`
	Name  string `json:"name"`
	State bool   `json:"state" gorm:"default:'true'"`

	Level1ID uint `json:"level_1_id"`
}
