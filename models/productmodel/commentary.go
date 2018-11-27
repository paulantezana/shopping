package productmodel

import "time"

type Commentary struct {
	ID           uint      `json:"id" gorm:"primary_key"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	UserName     string    `json:"user_name"`
	Email        string    `json:"email"`
	Country      string    `json:"country"`
	CommentTitle string    `json:"comment_title"`
	Commentary   string    `json:"commentary"`
	Points       string    `json:"points"`
	State        bool      `json:"state" gorm:"default:'true'"`

	ClientID  uint `json:"client_id"`
	ProductID uint `json:"product_id"`
}
