package models

import "time"

type Tax struct {
    ID    uint   `json:"id" gorm:"primary_key"`
    CreatedAt           time.Time `json:"created_at"`
    UpdatedAt           time.Time `json:"updated_at"`
    Name string `json:"name"`
    AcronymTax string `json:"acronym_tax"`
    TaxValue float32 `json:"tax_value"`
    Percentage bool `json:"percentage"`
    Default bool `json:"default"`
    InUse bool `json:"in_use"`
    State bool   `json:"state" gorm:"default:'true'"`
}
