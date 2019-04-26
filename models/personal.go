package models

import "time"

type Personal struct {
	ID             uint      `json:"id" gorm:"primary_key"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	DocNumber      string    `json:"doc_number"`
	Sex            string    `json:"sex"`
	Phone          string    `json:"phone"`
	Address        string    `json:"address"`
	BirthDate      time.Time `json:"birth_date"`
	State          uint8     `json:"state" gorm:"default:'1'"`
	TypeDocumentID uint      `json:"type_document_id"`

	// User
	User     string `json:"user"`
	Password string `json:"password"`
	Key  string `json:"key"`
	Avatar   string `json:"avatar"`
	Email    string `json:"email"`
    Freeze bool `json:"freeze"`
}
