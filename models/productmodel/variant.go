package productmodel

import "time"

type Variant struct {
	ID          uint      `json:"id" gorm:"primary_key"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	VariantName string    `json:"variant_name"`
	IsCombo     bool      `json:"is_combo"`
	State       bool      `json:"state" gorm:"default:'true'"`

	PresentationID uint `json:"presentation_id"`
}
