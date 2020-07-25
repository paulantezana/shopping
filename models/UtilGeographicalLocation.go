package models

// UtilGeographicalLocation --
type UtilGeographicalLocation struct {
	ID         uint   `json:"id" gorm:"primary_key"`
	Code       string `json:"code" gorm:"type:varchar(12)"`
	District   string `json:"district" gorm:"type:varchar(128)"`
	Province   string `json:"province" gorm:"type:varchar(128)"`
	Department string `json:"department" gorm:"type:varchar(128)"`
	State      bool   `json:"state" gorm:"default: true"`
}

// UtilGeographicalLocationShort --
type UtilGeographicalLocationShort struct {
	ID          uint   `json:"id" gorm:"-"`
	Code        string `json:"code" gorm:"-"`
	Description string `json:"description" gorm:"-"`
}
