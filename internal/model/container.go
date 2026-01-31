package model

type Container struct {
	PodID       uint       `json:"pod_id"`
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

func (c Container) GetModelName() string {
	return "Container"
}

func (c Container) GetBaseModel() BaseModel {
	return c.BaseModel
}

func (c Container) GetUniqueField() []string {
	return []string{"id"}
}

func (c Container) GetAllowedQueryFields() []string {
	return []string{}
}
