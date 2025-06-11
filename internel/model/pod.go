package model

import (
	"CBCTF/internel/config"
	"CBCTF/internel/i18n"
	"fmt"
)

type Pod struct {
	VictimID        uint            `json:"victim_id"`
	Victim          Victim          `json:"-"`
	Containers      []Container     `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	Name            string          `json:"name"`
	PodIP           string          `json:"pod_ip"`
	ExposedIP       string          `json:"exposed_ip"`
	PodPorts        Ports           `gorm:"type:json" json:"pod_ports"`
	ExposedPorts    Ports           `gorm:"type:json" json:"exposed_ports"`
	NetworkPolicies NetworkPolicies `gorm:"type:json" json:"network_policies"`
	Basic
}

func (p Pod) GetModelName() string {
	return "Pod"
}

func (p Pod) GetID() uint {
	return p.ID
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

func (p Pod) TrafficPath() string {
	return fmt.Sprintf("%s/traffics/victim-%d/pod-%d-%s.pcap", config.Env.Path, p.VictimID, p.ID, p.Name)
}

func (p Pod) RemoteAddr() []string {
	data := make([]string, 0)
	for _, port := range p.ExposedPorts {
		data = append(data, fmt.Sprintf("%s:%d", p.ExposedIP, port))
	}
	return data
}
