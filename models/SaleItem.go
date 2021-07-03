package models

type SaleItem struct {
	ID uint `json:"id" gorm:"primaryKey"`

	UtilUnitMeasureTypeId uint    `json:"util_unit_measure_type_id" gorm:"not null"`
	ProductCode           string  `json:"product_code" gorm:"not null"`
	Description           string  `json:"description" gorm:"not null"`
	Quantity              float64 `json:"quantity" gorm:"not null"`
	UnitValue             float64 `json:"unit_value" gorm:"not null"`
	UnitPrice             float64 `json:"unit_price" gorm:"not null"`

	Discount float64 `json:"discount" gorm:"default: 0"`
	Charge   float64 `json:"charge" gorm:"default: 0"`

	UtilAffectationIgvTypeId uint    `json:"util_affectation_igv_type_id"`
	TotalBaseIgv             float64 `json:"total_base_igv" gorm:"default: 0"`
	Igv                      float64 `json:"igv" gorm:"default: 0"`

	UtilSystemIscTypeId uint    `json:"util_system_isc_type_id" gorm:"default: 0"`
	TotalBaseIsc        float64 `json:"total_base_isc" gorm:"default: 0"`
	TaxIsc              float64 `json:"tax_isc" gorm:"default: 0"`
	Isc                 float64 `json:"isc" gorm:"default: 0"`

	TotalBaseOtherTaxed  float64 `json:"total_base_other_taxed" gorm:"default: 0"`
	PercentageOtherTaxed float64 `json:"percentage_other_taxed" gorm:"default: 0"`
	OtherTaxed           float64 `json:"other_taxed" gorm:"default: 0"`

	QuantityPlasticBag float64 `json:"quantity_plastic_bag" gorm:"default: 0"`
	PlasticBagTax      float64 `json:"plastic_bag_tax" gorm:"default: 0"`

	TotalValue float64 `json:"total_value" gorm:"default: 0"`
	Total      float64 `json:"total" gorm:"default: 0"`

	SaleId    uint `json:"sale_id" gorm:"not null"`
	ProductId uint `json:"product_id" gorm:"not null"`
}
