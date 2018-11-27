package productmodel

import (
	"time"
)

type Product struct {
	ID                      uint      `json:"id" gorm:"primary_key"`
	CreatedAt               time.Time `json:"created_at"`
	UpdatedAt               time.Time `json:"updated_at"`
	ProductCode             string    `json:"product_code"`
	BarCode                 string    `json:"bar_code"`
	ProductName             string    `json:"product_name"`
	LargeDescription        string    `json:"large_description"`
	ShortDescription        string    `json:"short_description"`
	Keywords                string    `json:"keywords"`
	PurchasePrice           float32   `json:"purchase_price"`
	InFrontPage             bool      `json:"in_front_page"`
	StarCategory            bool      `json:"star_category"`
	ControlWithoutStock     string    `json:"control_without_stock"`
	MinimumLimit            float32   `json:"minimum_limit"`
	MaximumLimit            float32   `json:"maximum_limit"`
	FractionAmount          bool      `json:"fraction_amount"`
	ProductBanner           string    `json:"product_banner"`
	ShowBanner              bool      `json:"show_banner"`
	UrlVideo                string    `json:"url_video"`
	ShowVideo               bool      `json:"show_video"`
	InUse                   bool      `json:"in_use"`
	SaleVariantWithoutStock bool      `json:"sale_variant_without_stock"`
	ShowWeb                 bool      `json:"show_web"`
	ShowWebPrice            bool      `json:"show_web_price"`
	State                   bool      `json:"state" gorm:"default:'true'"`

	UnitMeasureID uint `json:"unit_measure_id"`
	BrandID       uint `json:"brand_id"`
}
