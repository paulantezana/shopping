package models

type UtilCatAffectationIgvType struct {
	ID                uint   `json:"id" gorm:"primary_key"`
	Code              string `json:"code" gorm:"type:varchar(12)"`
	Description       string `json:"description" gorm:"type:varchar(128)"`
	Onerous           bool   `json:"onerous"`
	UtilTributeTypeId uint   `json:"util_tribute_type_id"`
	State             bool   `json:"state"`
}
