package model

import (
	"CBCTF/internal/traffic"
	"encoding/hex"
	"gorm.io/gorm"
	"time"
)

type Traffic struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	SrcIP       string         `json:"src_ip"`
	DstIP       string         `json:"dst_ip"`
	SrcPort     uint16         `json:"src_port"`
	DstPort     uint16         `json:"dst_port"`
	Payload     string         `json:"payload"`
	Type        string         `json:"type"`
	TeamID      uint           `json:"team"`
	ContestID   uint           `json:"contest"`
	ChallengeID string         `json:"challenge"`
	DockerID    uint           `json:"docker"`
	Path        string         `json:"path"`
	CreatedAt   time.Time      `json:"-"`
	UpdatedAt   time.Time      `json:"-"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func InitTraffic(conn traffic.Connection, docker Docker) Traffic {
	_, t := conn.ParsePayload()
	return Traffic{
		SrcIP:       conn.SrcIP,
		DstIP:       conn.DstIP,
		SrcPort:     conn.SrcPort,
		DstPort:     conn.DstPort,
		Payload:     hex.EncodeToString(conn.Payload),
		Type:        t,
		TeamID:      docker.TeamID,
		ContestID:   docker.ContestID,
		ChallengeID: docker.ChallengeID,
		DockerID:    docker.ID,
		Path:        docker.TrafficPath(),
	}
}
