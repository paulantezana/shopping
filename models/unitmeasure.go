package models

type UnitMeasure struct {
	ID     uint   `json:"id" gorm:"primary_key"`
	Name   string `json:"name"`
	Symbol string `json:"symbol"`
	State  bool   `json:"state" gorm:"default:'true'"`
}
