package models

type Warehouse struct {
	ID      uint   `json:"id" gorm:"primary_key"`
	Name    string `json:"name"`
	Primary bool   `json:"primary"`
	Address string `json:"address"`
	State   bool   `json:"state" gorm:"default:'true'"`

	BranchOfficeID       uint `json:"branch_office_id"`
	GeographicLocationID uint `json:"geographic_location_id"`
}
