package model

import (
	"fmt"
)

// Pod K8s Pod 实例
// BelongsTo Victim
// HasMany Container
type Pod struct {
	VictimID   uint        `json:"victim_id"`
	Victim     Victim      `json:"-"`
	Containers []Container `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	Name       string      `json:"name"`
	PodPorts   Exposes     `gorm:"type:json" json:"pod_ports"`
	Networks   Networks    `gorm:"type:json" json:"-"`
	BaseModel
}

func (p Pod) TableName() string {
	return "pods"
}

func (p Pod) ModelName() string {
	return "Pod"
}

func (p Pod) GetBaseModel() BaseModel {
	return p.BaseModel
}

func (p Pod) UniqueFields() []string {
	return []string{"id"}
}

func (p Pod) TrafficPcapPath() string {
	return fmt.Sprintf("%s/pod-%d.pcap", Victim{BaseModel: BaseModel{ID: p.VictimID}}.TrafficBasePath(), p.ID)
}

func (p Pod) QueryFields() []string {
	return []string{}
}
