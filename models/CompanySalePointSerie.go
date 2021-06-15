package models

import "time"

type CompanySalePointSerie struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	Serie       string `json:"serie" gorm:"type:varchar(4)"`
	Contingency bool   `json:"contingency"`

    CompanySalePointId uint `json:"company_sale_point_id"`
	UtilDocumentTypeId uint `json:"util_document_type_id"`

	CreatedAt     time.Time `json:"-"`
	UpdatedAt     time.Time `json:"-"`
	CreatedUserId uint      `json:"-"`
	UpdatedUserId uint      `json:"-"`
	State         bool      `json:"state" gorm:"default: true"`
}
