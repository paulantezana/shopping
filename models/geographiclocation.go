package models

type GeographicLocation struct {
	ID        uint `json:"id" gorm:"primary_key"`
	CountryID uint `json:"country_id"`
	Level1ID  uint `json:"level_1_id"`
	Level2ID  uint `json:"level_2_id"`
	Level3ID  uint `json:"level_3_id"`
}
