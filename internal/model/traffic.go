package model

type Traffic struct {
	VictimID uint   `gorm:"index;index:idx_traffics_victim_src_ip_active,priority:1,where:deleted_at IS NULL" json:"victim_id"`
	SrcIP    string `gorm:"index:idx_traffics_victim_src_ip_active,priority:2,where:deleted_at IS NULL" json:"src_ip"`
	DstIP    string `json:"dst_ip"`
	Type     string `json:"type"`
	Subtype  string `json:"subtype"`
	Count    uint   `json:"count"`
	Size     int    `json:"size"`
	BaseModel
}
