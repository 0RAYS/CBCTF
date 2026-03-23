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
