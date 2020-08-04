package models

import "time"

type CompanySerie struct {
	ID          uint   `json:"id" gorm:"primary_key"`
	Serie       string `json:"serie" gorm:"type:varchar(4)"`
	Contingency bool   `json:"contingency"`

	CompanyLocalId     uint `json:"company_local_id"`
	UtilDocumentTypeId uint `json:"util_document_type_id"`

	CreatedAt     time.Time `json:"-"`
	UpdatedAt     time.Time `json:"-"`
	CreatedUserId uint      `json:"-"`
	UpdatedUserId uint      `json:"-"`
	State         bool      `json:"state" gorm:"default: true"`
}
