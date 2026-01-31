package model

type Traffic struct {
	VictimID uint   `json:"victim_id"`
	SrcIP    string `json:"src_ip"`
	DstIP    string `json:"dst_ip"`
	Type     string `json:"type"`
	Subtype  string `json:"subtype"`
	Count    uint   `json:"count"`
	Size     int    `json:"size"`
	BaseModel
}

func (t Traffic) GetModelName() string {
	return "Traffic"
}

func (t Traffic) GetBaseModel() BaseModel {
	return t.BaseModel
}

func (t Traffic) GetUniqueField() []string {
	return []string{"id"}
}

func (t Traffic) GetAllowedQueryFields() []string {
	return []string{"id", "src_ip", "dst_ip", "type", "subtype"}
}
