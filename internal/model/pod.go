package model

import (
	"fmt"
)

type Pod struct {
	VictimID   uint        `json:"victim_id"`
	Containers []Container `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	Name       string      `json:"name"`
	PodPorts   Exposes     `gorm:"type:json" json:"pod_ports"`
	Networks   Networks    `gorm:"type:json" json:"-"`
	BaseModel
}

func (p Pod) GetModelName() string {
	return "Pod"
}

func (p Pod) GetBaseModel() BaseModel {
	return p.BaseModel
}

func (p Pod) GetUniqueField() []string {
	return []string{"id"}
}

func (p Pod) TrafficPcapPath() string {
	return fmt.Sprintf("%s/pod-%d.pcap", Victim{BaseModel: BaseModel{ID: p.VictimID}}.TrafficBasePath(), p.ID)
}

func (p Pod) GetAllowedQueryFields() []string {
	return []string{}
}
