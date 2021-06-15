package models

import "time"

type CompanyBankAccount struct {
	ID           uint   `json:"id" gorm:"primaryKey"`
	Account      string `json:"account"`
	Subsidiary   string `json:"subsidiary"`
	InterbankKey string `json:"interbank_key"`
	BankName     string `json:"bank_name"`
	CompanyId    uint   `json:"company_id"`

	CreatedAt     time.Time `json:"-"`
	UpdatedAt     time.Time `json:"-"`
	CreatedUserId uint      `json:"-"`
	UpdatedUserId uint      `json:"-"`
	State         bool      `json:"state" gorm:"default: true"`
}
