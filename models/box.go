package models

type Box struct {
	ID    uint   `json:"id" gorm:"primary_key"`
	Name  string `json:"name"`
	State bool   `json:"state" gorm:"default:'true'"`

	BranchOffice uint `json:"branch_office"`
}
