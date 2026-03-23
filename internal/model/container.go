package model

// Container K8s 容器
// BelongsTo Pod
type Container struct {
	PodID       uint       `json:"pod_id"`
	Pod         Pod        `json:"-"`
	Name        string     `json:"name"`
	Image       string     `json:"image"`
	CPU         float32    `json:"cpu"`
	Memory      int64      `json:"memory"`
	WorkingDir  string     `gorm:"default:null" json:"working_dir"`
	Command     StringList `gorm:"default:null;type:json" json:"command"`
	Environment StringMap  `gorm:"default:null;type:json" json:"environment"`
	EnvFlags    StringMap  `gorm:"type:json" json:"env_flags"`
	VolumeFlags StringMap  `gorm:"type:json" json:"volume_flags"`
	Exposes     Exposes    `gorm:"type:json" json:"exposes"`
	BaseModel
}
