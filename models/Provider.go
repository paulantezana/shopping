package models

import "time"

type Provider struct {
	ID             uint   `json:"id" gorm:"primary_key"`
	DocumentNumber string `json:"document_number" gorm:"type:varchar(32)"`
	SocialReason   string `json:"social_reason" gorm:"type:varchar(255)"`
	Phone          string `json:"phone" gorm:"type:varchar(32)"`
	Address        string `json:"address" gorm:"type:varchar(255)"`
	Email          string `json:"email" gorm:"type:varchar(64)"`
	Observation    string `json:"observation" gorm:"type:varchar(255)"`

	CompanyId                  uint `json:"company_id"`
	UtilIdentityDocumentTypeId uint `json:"util_identity_document_type_id"`

	CreatedAt     time.Time `json:"-"`
	UpdatedAt     time.Time `json:"-"`
	CreatedUserId uint      `json:"-"`
	UpdatedUserId uint      `json:"-"`
	State         bool      `json:"state" gorm:"default: true"`
}
