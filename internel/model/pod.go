package model

import (
	"CBCTF/internel/config"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type Pod struct {
	ID                uint            `gorm:"primaryKey" json:"id"`
	VictimID          uint            `json:"victim_id"`
	Victim            Victim          `json:"-"`
	Containers        []Container     `json:"-"`
	Traffics          []Traffic       `json:"-"`
	Name              string          `json:"name"`
	ExposeIP          string          `json:"ip"`
	PodIP             string          `json:"pod_ip"`
	ServiceName       string          `json:"service"`
	NetworkPolicyName string          `json:"network_policy"`
	ExposePorts       Ports           `gorm:"type:json" json:"exposes"`
	NetworkPolicies   NetworkPolicies `json:"network_policies"`
	CreatedAt         time.Time       `json:"-"`
	UpdatedAt         time.Time       `json:"-"`
	DeletedAt         gorm.DeletedAt  `gorm:"index" json:"-"`
	Version           uint            `gorm:"default:1" json:"-"`
}

func (p Pod) TrafficPath() string {
	return fmt.Sprintf("%s/traffics/%d/%s/%d.pcap", config.Env.Path, p.VictimID, p.Name, p.ID)
}

func (p Pod) RemoteAddr() []string {
	data := make([]string, 0)
	for _, port := range p.ExposePorts {
		data = append(data, fmt.Sprintf("%s:%d", p.ExposeIP, port))
	}
	return data
}
