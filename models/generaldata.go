package models

type GeneralData struct {
	ID           uint   `json:"id" gorm:"primary_key"`
	BusinessName string `json:"business_name"`
	Ruc          string `json:"ruc"`
	Address      string `json:"address"`
	CompanyLogo  string `json:"company_logo"`
	Email        string `json:"email"`
	BankAccount  string `json:"bank_account"`
	State        bool   `json:"state" gorm:"default:'true'"`

	GeographicLocationID uint `json:"geographic_location_id"`
}
