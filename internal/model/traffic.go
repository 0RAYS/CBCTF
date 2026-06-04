package model

// Traffic 每个 victim 一条，IPs 存储靶机流量中涉及的所有 IP（普通 pod + frpc proxy protocol）。
type Traffic struct {
	VictimID uint       `gorm:"uniqueIndex:idx_traffics_victim_id_active,where:deleted_at IS NULL;index" json:"victim_id"`
	IPs      StringList `gorm:"type:jsonb;default:'[]'" json:"ips"`
	BaseModel
}
