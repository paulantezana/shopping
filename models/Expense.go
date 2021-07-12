package models

import "time"

type Expense struct {
	ID                 uint      `json:"id" gorm:"primaryKey"`
	DateOfIssue        time.Time `json:"date_off_issue"`
	Description        string    `json:"description"`
	Amount             float32   `json:"amount"`
	CurrencyId         uint      `json:"currency_id"`
	ExpenseTypeId      uint      `json:"expense_type_id"`
	UserId             uint      `json:"user_id"`
	CompanySalePointId uint      `json:"company_sale_point_id"`

	CreatedAt     time.Time `json:"-"`
	UpdatedAt     time.Time `json:"-"`
	CreatedUserId uint      `json:"-"`
	UpdatedUserId uint      `json:"-"`
	State         bool      `json:"state" gorm:"default: true"`
}
