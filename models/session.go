package models

import "time"

type Session struct {
	ID           uint      `json:"id" gorm:"primary_key"`
	IpAddress    string    `json:"ip_address"`
	UserAgent    string    `json:"user_agent"`
	LastActivity time.Time `json:"last_activity"`
	UserData     string    `json:"user_data"`
}
