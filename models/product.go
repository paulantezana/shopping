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
	UnitMeasureId uint `json:"unit_measure_id"`

	CreatedAt     time.Time `json:"-"`
	UpdatedAt     time.Time `json:"-"`
	CreatedUserId uint      `json:"-"`
	UpdatedUserId uint      `json:"-"`
	State         bool      `json:"state" gorm:"default: true"`
}
