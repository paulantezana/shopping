package models

import "time"

type Product struct {
	ID uint `json:"id" gorm:"primary_key"`

	Url             string  `json:"url" gorm:"not null; type:varchar(255)"`
	Title           string  `json:"title" gorm:"not null; type:varchar(255)"`
	Description     string  `json:"description" gorm:"not null; type:varchar(255)"`
	LongDescription string  `json:"long_description"`
	Barcode         string  `json:"barcode" gorm:"not null; type:varchar(32); unique_index:idx_product_key"`
	IsService       bool    `json:"is_service"`
	Location        string  `json:"location"`
	StockMin        float32 `json:"stock_min" gorm:"default: 0.00"`
	StockMax        float32 `json:"stock_max" gorm:"default: 0.00"`

	InternalUse   bool    `json:"internal_use"`
	Favourite     bool    `json:"favourite"`
	PurchasePrice float64 `json:"purchase_price"`

	Lot    bool    `json:"lot"`                         // Lote
	Bulk   bool    `json:"bulk"`                        // Granel
	Recipe bool    `json:"recipe"`                      // Receta medica
	Weight float32 `json:"weight" gorm:"default: 0.00"` // Peso

	SalePrice1 float64 `json:"sale_price_1" gorm:"default: 0.00"`
	SalePrice2 float64 `json:"sale_price_2" gorm:"default: 0.00"`
	SalePrice3 float64 `json:"sale_price_3" gorm:"default: 0.00"`
	SalePrice4 float64 `json:"sale_price_4" gorm:"default: 0.00"`

	WholeSale1 float64 `json:"whole_sale_1" gorm:"default: 0.00"`
	WholeSale2 float64 `json:"whole_sale_2" gorm:"default: 0.00"`
	WholeSale3 float64 `json:"whole_sale_3" gorm:"default: 0.00"`
	WholeSale4 float64 `json:"whole_sale_4" gorm:"default: 0.00"`

	PurchaseUtilUnitMeasureTypeId uint    `json:"purchase_util_unit_measure_type_id"`
	SaleUtilUnitMeasureTypeId     uint    `json:"sale_util_unit_measure_type_id"`
	Factor                        float32 `json:"factor"`

	CompanyId  uint `json:"company_id" gorm:"not null; unique_index:idx_product_key"`
	CategoryId uint `json:"category_id"`

	CreatedAt     time.Time `json:"-"`
	UpdatedAt     time.Time `json:"-"`
	CreatedUserId uint      `json:"-"`
	UpdatedUserId uint      `json:"-"`
	State         bool      `json:"state" gorm:"default: true"`

	// AUX
	Stock                          float32 `json:"stock" gorm:"-"`
	PurchaseUnitMeasureDescription string  `json:"purchase_unit_measure_description" gorm:"-"`
	SearchText                     string  `json:"search_text" gorm:"-"`
}
