package model

import (
	"CBCTF/internal/i18n"
	"fmt"
)

type Pod struct {
	VictimID   uint        `json:"victim_id"`
	Victim     Victim      `json:"-"`
	Containers []Container `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	Traffics   []Traffic   `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	Name       string      `json:"name"`
	PodPorts   Exposes     `gorm:"type:json" json:"pod_ports"`
	Networks   Networks    `gorm:"type:json" json:"-"`
	BasicModel
}

func (p Pod) GetModelName() string {
	return "Pod"
}

func (p Pod) GetVersion() uint {
	return p.Version
}

func (p Pod) CreateErrorString() string {
	return i18n.CreatePodError
}

func (p Pod) DeleteErrorString() string {
	return i18n.DeletePodError
}

func (p Pod) GetErrorString() string {
	return i18n.GetPodError
}

func (p Pod) NotFoundErrorString() string {
	return i18n.PodNotFound
}

func (p Pod) UpdateErrorString() string {
	return i18n.UpdatePodError
}

func (p Pod) GetUniqueKey() []string {
	return []string{"id"}
}

func (p Pod) GetForeignKeys() []string {
	return []string{"id", "victim_id"}
}

func (p Pod) TrafficPcapPath() string {
	return fmt.Sprintf("%s/pod-%d.pcap", Victim{BasicModel: BasicModel{ID: p.VictimID}}.TrafficBasePath(), p.ID)
}
