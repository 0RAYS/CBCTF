package model

import (
	"CBCTF/internel/config"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type Pod struct {
	ID              uint            `gorm:"primaryKey" json:"id"`
	VictimID        uint            `json:"victim_id"`
	Victim          Victim          `json:"-"`
	Containers      []Container     `json:"-"`
	Traffics        []Traffic       `json:"-"`
	Name            string          `json:"name"`
	PodIP           string          `json:"pod_ip"`
	ExposedIP       string          `json:"exposed_ip"`
	PodPorts        Ports           `gorm:"type:json" json:"pod_ports"`
	ExposedPorts    Ports           `gorm:"type:json" json:"exposed_ports"`
	NetworkPolicies NetworkPolicies `gorm:"type:json" json:"network_policies"`
	CreatedAt       time.Time       `json:"-"`
	UpdatedAt       time.Time       `json:"-"`
	DeletedAt       gorm.DeletedAt  `gorm:"index" json:"-"`
	Version         uint            `gorm:"default:1" json:"-"`
}

func (p Pod) TrafficPath() string {
	return fmt.Sprintf("%s/traffics/victim-%d/%s/pod-%d.pcap", config.Env.Path, p.VictimID, p.Name, p.ID)
}

func (p Pod) RemoteAddr() []string {
	data := make([]string, 0)
	for _, port := range p.ExposedPorts {
		data = append(data, fmt.Sprintf("%s:%d", p.ExposedIP, port))
	}
	return data
}
