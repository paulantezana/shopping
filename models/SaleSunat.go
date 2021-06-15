package models

type SaleSunat struct {
	ID uint `json:"id" gorm:"primaryKey"`

	PdfUrl          string `json:"pdf_url" gorm:"default: ''"`
	XmlUrl          string `json:"xml_url" gorm:"default: ''"`
	CdrUrl          string `json:"cdr_url" gorm:"default: ''"`
	Send            bool   `json:"send" gorm:"default: false"`
	ResponseCode    string `json:"response_code" gorm:"default: ''"`
	ResponseMessage string `json:"response_message" gorm:"default: ''"`
	OtherMessage    string `json:"other_message" gorm:"default: ''"`

	SaleId uint `json:"sale_id" gorm:"unique; not null"`
}
