package models

// UtilDocumentType --
type UtilDocumentType struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	Code        string `json:"code" gorm:"type:varchar(12)"`
	Description string `json:"description" gorm:"type:varchar(128)"`
	Sunat       bool   `json:"sunat" gorm:"default: true"`
	State       bool   `json:"state" gorm:"default: true"`
}
