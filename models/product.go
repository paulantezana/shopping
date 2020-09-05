package models

import "time"

type Product struct {
	ID uint `json:"id" gorm:"primary_key"`

	Url             string `json:"url"`
	Title           string `json:"title"`
	Description     string `json:"description"`
	LongDescription string `json:"long_description"`

	CompanyId     uint `json:"company_id"`
	CategoryId    uint `json:"category_id"`
	UtilUnitMeasureType uint `json:"util_unit_measure_type"`

	CreatedAt     time.Time `json:"-"`
	UpdatedAt     time.Time `json:"-"`
	CreatedUserId uint      `json:"-"`
	UpdatedUserId uint      `json:"-"`
	State         bool      `json:"state" gorm:"default: true"`
}
