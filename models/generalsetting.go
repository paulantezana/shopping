package models

type GeneralSetting struct {
	ID                uint    `json:"id" gorm:"primary_key"`
	DecimalNumbers    float32 `json:"decimal_numbers"`
	ImageDefault      string  `json:"image_default"`
	PercentageUtility float32 `json:"percentage_utility"`
	ArchingMarker     bool    `json:"arching_marker"`
	LogoImpression    string  `json:"logo_impression"`
	State             bool    `json:"state" gorm:"default:'true'"`

	GeneralDataID uint `json:"general_data_id"`
}
