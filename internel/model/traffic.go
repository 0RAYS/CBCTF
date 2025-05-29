package model

import (
	"time"
)

type Traffic struct {
	VictimID uint      `json:"victim_id"`
	Victim   Victim    `json:"-"`
	PodID    uint      `json:"pod_id"`
	Pod      Pod       `json:"-"`
	SrcIP    string    `json:"src_ip"`
	DstIP    string    `json:"dst_ip"`
	SrcPort  uint16    `json:"src_port"`
	DstPort  uint16    `json:"dst_port"`
	Payload  string    `json:"payload"`
	Time     time.Time `json:"time"`
	Type     string    `json:"type"`
	Path     string    `json:"path"`
	BaseModel
}
