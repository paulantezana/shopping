package models

// UtilTributeType --
type UtilTributeType struct {
	ID                uint   `json:"id" gorm:"primaryKey"`
	Code              string `json:"code" gorm:"type:varchar(12)"`
	Description       string `json:"description" gorm:"type:varchar(128)"`
	InternationalCode string `json:"international_code" gorm:"type:varchar(12)"`
	Name              string `json:"name" gorm:"type:varchar(12)"`
	State             bool   `json:"state" gorm:"default: true"`
}
