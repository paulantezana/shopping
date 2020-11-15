package models

import "time"

type Product struct {
	ID uint `json:"id" gorm:"primary_key"`

	Url             string  `json:"url"`
	Title           string  `json:"title"`
	Description     string  `json:"description"`
	LongDescription string  `json:"long_description"`
	Barcode         string  `json:"barcode"`
	BarcodeAux      string  `json:"barcode_aux"`
	IsService       bool    `json:"is_service"`
	Location        string  `json:"location"`
	StockMin        float32 `json:"stock_min"`
	StockMax        float32 `json:"stock_max"`
	Weight          float32 `json:"weight"`
	Recipe          bool    `json:"recipe"`
	Bulk            bool    `json:"bulk"`
	InternalUse     bool    `json:"internal_use"`
	Favourite       bool    `json:"favourite"`
	PurchasePrice   float64 `json:"purchase_price"`
	SalePrice1      float64 `json:"sale_price_1"`
	SalePrice2      float64 `json:"sale_price_2"`
	SalePrice3      float64 `json:"sale_price_3"`
	WholeSale1      float64 `json:"whole_sale_1"`
	WholeSale2      float64 `json:"whole_sale_2"`
	WholeSale3      float64 `json:"whole_sale_3"`

	CompanyId                     uint    `json:"company_id"`
	CategoryId                    uint    `json:"category_id"`
	BrandId                       uint    `json:"brand_id"`
	PatternId                     uint    `json:"pattern_id"`
	PurchaseUtilUnitMeasureTypeId uint    `json:"purchase_util_unit_measure_type_id"`
	SaleUtilUnitMeasureTypeId     uint    `json:"sale_util_unit_measure_type_id"`
	Factor                        float32 `json:"factor"`

	CreatedAt     time.Time `json:"-"`
	UpdatedAt     time.Time `json:"-"`
	CreatedUserId uint      `json:"-"`
	UpdatedUserId uint      `json:"-"`
	State         bool      `json:"state" gorm:"default: true"`

	// AUX
	Stock float32 `json:"stock"  gorm:"-"`
}
