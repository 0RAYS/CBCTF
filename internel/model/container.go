package model

import (
	"gorm.io/gorm"
	"time"
)

type Container struct {
	ID                uint           `gorm:"primaryKey" json:"id"`
	UsageID           uint           `json:"usage_id"`
	Usage             Usage          `json:"-"`
	TeamID            uint           `json:"team_id"`
	Team              Team           `json:"-"`
	UserID            uint           `json:"user_id"`
	User              User           `json:"-"`
	Exposes           Exposes        `gorm:"type:json" json:"exposes"`
	Start             time.Time      `json:"start"`
	Duration          time.Duration  `json:"duration"`
	PodName           string         `json:"pod"`
	ContainerName     string         `json:"container"`
	ServiceName       string         `json:"service"`
	NetworkPolicyName string         `json:"network_policy"`
	Flags             Strings        `gorm:"type:json" json:"flags"`
	Traffics          []Traffic      `json:"-"`
	CreatedAt         time.Time      `json:"-"`
	UpdatedAt         time.Time      `json:"-"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`
	Version           uint           `gorm:"default:1" json:"-"`
}
