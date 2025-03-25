package model

import (
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
	ContainerID uint           `json:"container_id"`
	Container   Container      `json:"-"`
	Path        string         `json:"path"`
	CreatedAt   time.Time      `json:"-"`
	UpdatedAt   time.Time      `json:"-"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
