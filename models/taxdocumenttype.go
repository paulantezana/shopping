package models

type TaxDocumentType struct {
    ID    uint   `json:"id" gorm:"primary_key"`

    State bool   `json:"state" gorm:"default:'true'"`
} 
