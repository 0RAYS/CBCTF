package model

import "CBCTF/internel/i18n"

type Traffic struct {
	VictimID uint   `json:"victim_id"`
	Victim   Victim `json:"-"`
	PodID    uint   `json:"pod_id"`
	Pod      Pod    `json:"-"`
	SrcIP    string `json:"src_ip"`
	DstIP    string `json:"dst_ip"`
	SrcPort  uint16 `json:"src_port"`
	DstPort  uint16 `json:"dst_port"`
	Type     string `json:"type"`
	Count    uint   `json:"count"`
	BasicModel
}

func (t Traffic) GetModelName() string {
	return "Traffic"
}

func (t Traffic) GetID() uint {
	return t.ID
}

func (t Traffic) GetVersion() uint {
	return t.Version
}

func (t Traffic) CreateErrorString() string {
	return i18n.CreateTrafficError
}

func (t Traffic) DeleteErrorString() string {
	return i18n.DeleteTrafficError
}

func (t Traffic) GetErrorString() string {
	return i18n.GetTrafficError
}

func (t Traffic) NotFoundErrorString() string {
	return i18n.TrafficNotFound
}

func (t Traffic) UpdateErrorString() string {
	return i18n.UpdateTrafficError
}

func (t Traffic) GetUniqueKey() []string {
	return []string{"id"}
}
