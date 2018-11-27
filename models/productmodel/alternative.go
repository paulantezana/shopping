package productmodel

type Alternative struct {
	ID          uint   `json:"id" gorm:"primary_key"`
	Description string `json:"description"`
	Selected    bool   `json:"selected"`
	Position    uint   `json:"position"`
	State       bool   `json:"state" gorm:"default:'true'"`

	VariantID uint `json:"variant_id"`
}
