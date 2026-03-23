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
	Command     StringList `gorm:"default:null;type:jsonb" json:"command"`
	Environment StringMap  `gorm:"default:null;type:jsonb" json:"environment"`
	EnvFlags    StringMap  `gorm:"type:jsonb" json:"env_flags"`
	VolumeFlags StringMap  `gorm:"type:jsonb" json:"volume_flags"`
	Exposes     Exposes    `gorm:"type:jsonb" json:"exposes"`
	BaseModel
}
