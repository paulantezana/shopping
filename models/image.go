package models

type Image struct {
	ID             uint   `json:"id" gorm:"primary_key"`
	LocationImage  string `json:"location_image"`
	NameImage      string `json:"name_image"`
	ExtensionImage string `json:"extension_image"`
	CaptionImage   string `json:"caption_image"`
	ItIsCover      bool   `json:"it_is_cover"` // En portada
	State          bool   `json:"state" gorm:"default:'true'"`

	ProductID uint `json:"product_id"`
}
