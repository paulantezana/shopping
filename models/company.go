package models

type Company struct {
	ID           uint   `json:"id" gorm:"primary_key"`
    CompanyName string `json:"company_name"`
	Ruc          string `json:"ruc"`
	Address      string `json:"address"`
	CompanyLogo  string `json:"company_logo"`
	Email        string `json:"email"`
	BankAccount  string `json:"bank_account"`
	State        bool   `json:"state" gorm:"default:'true'"`
}
