package models

import "time"

type CompanySalePointDocumentAuth struct {
	ID uint `json:"id" gorm:"primary_key"`

	CompanySalePoint   uint `json:"company_sale_point"`
	UtilDocumentTypeId uint `json:"util_document_type_id"`

	CreatedAt     time.Time `json:"-"`
	UpdatedAt     time.Time `json:"-"`
	CreatedUserId uint      `json:"-"`
	UpdatedUserId uint      `json:"-"`
	State         bool      `json:"state" gorm:"default: true"`
}
