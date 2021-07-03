package models

// UtilAffectationIgvType --
type UtilAffectationIgvType struct {
	ID                uint   `json:"id" gorm:"primaryKey"`
	Code              string `json:"code" gorm:"type:varchar(12)"`
	NuCode            string `json:"nu_code" gorm:"type:varchar(5)"`
	Description       string `json:"description" gorm:"type:varchar(128)"`
	Onerous           bool   `json:"onerous"`
	UtilTributeTypeId uint   `json:"util_tribute_type_id"`
	State             bool   `json:"state" gorm:"default: true"`
}
