package productmodel

import (
	"time"
)

type Presentation struct {
	ID                  uint      `json:"id" gorm:"primary_key"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
	Code                string    `json:"code"`
	Barcode             string    `json:"barcode"`
	PresentationName    string    `json:"presentation_name"`
	Description         string    `json:"description"`
	PresentationSymbol  string    `json:"presentation_symbol"`
	PurchasePrice       float32   `json:"purchase_price"`
	UnitAmount          float32   `json:"unit_amount"`
	DefaultPresentation bool      `json:"default_presentation"`
	State               bool      `json:"state" gorm:"default:'true'"`

	ProductID          uint `json:"product_id"`
	PresentationBaseID uint `json:"presentation_base_id"`
}
