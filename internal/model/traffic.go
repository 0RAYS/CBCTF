package model

import "CBCTF/internal/traffic"

// Traffic 每个 victim 一条，IPs 存储靶机流量中涉及的所有 IP（普通 pod + frpc proxy protocol）。
type Traffic struct {
	VictimID uint                    `gorm:"uniqueIndex;index" json:"victim_id"`
	IPs      StringList              `gorm:"type:jsonb;default:'[]'" json:"ips"`
	Accesses []traffic.TrafficAccess `gorm:"serializer:json;type:jsonb;default:'[]'" json:"accesses"`
	Analysis *traffic.AnalysisReport `gorm:"serializer:json;type:jsonb" json:"analysis"`
	BaseModel
}
