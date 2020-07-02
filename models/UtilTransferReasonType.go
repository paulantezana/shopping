package models

// UtilTransferReasonType --
type UtilTransferReasonType struct {
	ID          uint   `json:"id" gorm:"primary_key"`
	Code        string `json:"code" gorm:"type:varchar(12)"`
	Description string `json:"description" gorm:"type:varchar(128)"`
	State       bool   `json:"state" gorm:"default: true"`
}
