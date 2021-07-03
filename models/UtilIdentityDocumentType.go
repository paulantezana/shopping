package models

// UtilIdentityDocumentType --
type UtilIdentityDocumentType struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	Code        string `json:"code" gorm:"type:varchar(12)"`
	NuCode      string `json:"nu_code" gorm:"type:varchar(5)"`
	Description string `json:"description" gorm:"type:varchar(128)"`
	State       bool   `json:"state" gorm:"default: true"`
}
