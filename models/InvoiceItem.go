package models

type InvoiceItem struct {
	ID uint `json:"id" gorm:"primary_key"`

	UnitMeasure string  `json:"unit_measure" gorm:"not null"`
	ProductCode string  `json:"product_code" gorm:"not null"`
	Description string  `json:"description" gorm:"not null"`
	Quantity    float64 `json:"quantity" gorm:"not null"`
	UnitValue   float64 `json:"unit_value" gorm:"not null"`
	UnitPrice   float64 `json:"unit_price" gorm:"not null"`

	Discount float64 `json:"discount" gorm:"default: 0.00"`
	Charge   float64 `json:"charge" gorm:"default: 0.00"`

	UtilCatAffectationIgvTypeId uint    `json:"util_cat_affectation_igv_type_id"`
	TotalBaseIgv                float64 `json:"total_base_igv" gorm:"default: 0.00"`
	Igv                         float64 `json:"igv" gorm:"default: 0.00"`

	UtilSystemIscTypeId uint    `json:"util_system_isc_type_id" gorm:"default: 0.00"`
	TotalBaseIsc        float64 `json:"total_base_isc" gorm:"default: 0.00"`
	TaxIsc              float64 `json:"tax_isc" gorm:"default: 0.00"`
	Isc                 float64 `json:"isc" gorm:"default: 0.00"`

	TotalBaseOtherTaxed  float64 `json:"total_base_other_taxed" gorm:"default: 0.00"`
	PercentageOtherTaxed float64 `json:"percentage_other_taxed" gorm:"default: 0.00"`
	OtherTaxed           float64 `json:"other_taxed" gorm:"default: 0.00"`

	QuantityPlasticBag float64 `json:"quantity_plastic_bag" gorm:"default: 0.00"`
	PlasticBagTax      float64 `json:"plastic_bag_tax" gorm:"default: 0.00"`

	TotalValue float64 `json:"total_value" gorm:"default: 0.00"`
	Total      float64 `json:"total" gorm:"default: 0.00"`

	InvoiceId uint `json:"invoice_id" gorm:"not null"`
}
