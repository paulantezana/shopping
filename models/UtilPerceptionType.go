package models

// UtilPerceptionType --
type UtilPerceptionType struct {
	ID          uint    `json:"id" gorm:"primaryKey"`
	Code        string  `json:"code" gorm:"type:varchar(12)"`
	Description string  `json:"description" gorm:"type:varchar(128)"`
	Percentage  float32 `json:"percentage"`
	State       bool    `json:"state" gorm:"default: true"`
}
