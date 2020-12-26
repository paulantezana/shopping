package models

import "time"

type ProductWareHouse struct {
	ID          uint    `json:"id" gorm:"primary_key"`
	Stock       float64 `json:"stock"`
	ProductId   uint    `json:"product_id"`
	WareHouseId uint    `json:"ware_house_id"`

	CreatedAt     time.Time `json:"-"`
	UpdatedAt     time.Time `json:"-"`
	CreatedUserId uint      `json:"-"`
	UpdatedUserId uint      `json:"-"`
	State         bool      `json:"state" gorm:"default: true"`
}
