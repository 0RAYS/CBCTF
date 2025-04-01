package model

import (
	"CBCTF/internel/config"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type Container struct {
	ID                uint            `gorm:"primaryKey" json:"id"`
	UsageID           uint            `json:"usage_id"`
	Usage             Usage           `json:"-"`
	TeamID            uint            `json:"team_id"`
	Team              Team            `json:"-"`
	UserID            uint            `json:"user_id"`
	User              User            `json:"-"`
	IP                string          `json:"ip"`
	Exposes           Exposes         `gorm:"type:json" json:"exposes"`
	Start             time.Time       `json:"start"`
	Duration          time.Duration   `json:"duration"`
	Image             string          `json:"image"`
	PodName           string          `json:"pod"`
	ContainerName     string          `json:"container"`
	ServiceName       string          `json:"service"`
	NetworkPolicyName string          `json:"network_policy"`
	NetworkPolicies   NetworkPolicies `json:"network_policies"`
	Flags             Strings         `gorm:"type:json" json:"flags"`
	Traffics          []Traffic       `json:"-"`
	CreatedAt         time.Time       `json:"-"`
	UpdatedAt         time.Time       `json:"-"`
	DeletedAt         gorm.DeletedAt  `gorm:"index" json:"-"`
	Version           uint            `gorm:"default:1" json:"-"`
}

func (c Container) TrafficPath() string {
	return fmt.Sprintf("%s/traffics/%s/%d/%d.pcap", config.Env.Path, c.PodName, c.TeamID, c.ID)
}

// RemoteAddr 返回远程地址
func (c Container) RemoteAddr() []string {
	data := make([]string, 0)
	for _, port := range c.Exposes {
		data = append(data, fmt.Sprintf("%s:%d", c.IP, port))
	}
	return data
}

// Remaining 返回剩余时间
func (c Container) Remaining() time.Duration {
	return c.Start.Add(c.Duration).Sub(time.Now())
}
