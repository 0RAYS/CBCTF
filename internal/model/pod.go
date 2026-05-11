package model

import "path/filepath"

// Pod K8s Pod 实例
// BelongsTo Victim
type Pod struct {
	VictimID uint    `json:"victim_id"`
	Victim   Victim  `json:"-"`
	Name     string  `json:"name"`
	Spec     PodSpec `gorm:"type:jsonb" json:"-"`
	BaseModel
}

func (p Pod) TrafficPcapPath() string {
	return filepath.Join(Victim{BaseModel: BaseModel{ID: p.VictimID}}.TrafficBasePath(), "pod-"+p.Name+".pcap")
}
