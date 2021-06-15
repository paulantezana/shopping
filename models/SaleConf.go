package models

type SaleConf struct {
	ID             uint   `json:"id" gorm:"primaryKey"`
	DocumentSize   string `json:"document_size"`
	SaleOutOfStock bool   `json:"sale_out_of_stock"`
	ItemsWithTaxes bool   `json:"items_with_taxes"`

	CompanyId uint `json:"company_id"`
}
