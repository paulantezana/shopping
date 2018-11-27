package models

type Level3 struct {
	ID    uint   `json:"id" gorm:"primary_key"`
	Name  string `json:"name"`
	State bool   `json:"state" gorm:"default:'true'"`

	Level2ID uint `json:"level_2_id"`
}
