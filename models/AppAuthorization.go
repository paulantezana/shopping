package models

// AppAuthorization --
type AppAuthorization struct {
	ID          uint   `json:"id" gorm:"primary_key"`
	Key         string `json:"key" gorm:"unique; not null; type:varchar(64)"`
	Type        uint   `json:"type" gorm:"default: 0"`                                            // 0 = menu, 1 = action
	Action      string `json:"action" gorm:"not null; type:varchar(128)"`
	Description string `json:"description" gorm:"not null; type:varchar(255)"`
	ParentId    uint   `json:"parent_id" gorm:"default: 0"`
	State       bool   `json:"state" gorm:"default: true"`
}
