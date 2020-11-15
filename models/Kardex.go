package models

import "time"

type Kardex struct {
	ID                  uint      `json:"id" gorm:"primary_key"`
	DateOfIssue         time.Time `json:"date_of_issue" gorm:"not null"`
	Quantity            float64   `json:"quantity" gorm:"not null"`
	UnitPrice           float64   `json:"unit_price" gorm:"not null"`
	Total               float64   `json:"total" gorm:"default: 0.00"`
	Stock               float64   `json:"stock"`
	Origin              string    `json:"origin"`
	Destination         string    `json:"destination"`
	Description         string    `json:"description"`
	DocumentDescription string    `json:"document_description"`

	UserId             uint `json:"user_id"`
	ProductId          uint `json:"product_id"`
	CompanyWareHouseId uint `json:"company_ware_house_id"`

	CreatedAt     time.Time `json:"-"`
	UpdatedAt     time.Time `json:"-"`
	CreatedUserId uint      `json:"-"`
	UpdatedUserId uint      `json:"-"`
	State         bool      `json:"state" gorm:"default: true"`
}
