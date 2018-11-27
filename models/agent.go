package models

type Agent struct {
	ID          uint `json:"id" gorm:"primary_key"`
	TypeAgentID uint `json:"type_agent_id"`
	ClientID    uint `json:"client_id"`
}
