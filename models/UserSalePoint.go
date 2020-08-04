package models

import "time"

type UserSalePoint struct {
    CreatedAt     time.Time `json:"-"`
    UpdatedAt     time.Time `json:"-"`
    CreatedUserId uint      `json:"-"`
    UpdatedUserId uint      `json:"-"`
    State         bool      `json:"state" gorm:"default: true"`
}
