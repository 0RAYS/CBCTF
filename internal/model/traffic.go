package model

import "CBCTF/internal/i18n"

type Traffic struct {
	VictimID uint   `json:"victim_id"`
	Victim   Victim `json:"-"`
	SrcIP    string `json:"src_ip"`
	DstIP    string `json:"dst_ip"`
	Type     string `json:"type"`
	Subtype  string `json:"subtype"`
	Count    uint   `json:"count"`
	Size     int    `json:"size"`
	BasicModel
}

func (t Traffic) GetModelName() string {
	return "Traffic"
}

func (t Traffic) GetBasicModel() BasicModel {
	return t.BasicModel
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
