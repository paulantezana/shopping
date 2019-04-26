package models

type DocumentType struct {
	ID             uint   `json:"id" gorm:"primary_key"`
	Name           string `json:"name"`
	LabelName      string `json:"label_name"`
	Area           string `json:"area"`
	Description    string `json:"description"`
	Voucher        uint   `json:"voucher"`
	DocumentFormat string `json:"document_format"`
	ResizeModel    uint   `json:"resize_model"`
	DetailEdge     uint   `json:"detail_edge"`
	ClientType     uint   `json:"client_type"`
	State          bool   `json:"state" gorm:"default:'true'"`
}
