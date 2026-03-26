package model

import (
	"fmt"
)

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
	return fmt.Sprintf("%s/pod-%s.pcap", Victim{BaseModel: BaseModel{ID: p.VictimID}}.TrafficBasePath(), p.Name)
}
