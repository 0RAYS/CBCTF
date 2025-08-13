package model

import "CBCTF/internal/i18n"

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
	BasicModel
}

func (c Container) GetModelName() string {
	return "Container"
}

func (c Container) GetVersion() uint {
	return c.Version
}

func (c Container) GetBasicModel() BasicModel {
	return c.BasicModel
}

func (c Container) CreateErrorString() string {
	return i18n.CreateContainerError
}

func (c Container) DeleteErrorString() string {
	return i18n.DeleteContainerError
}

func (c Container) GetErrorString() string {
	return i18n.GetContainerError
}

func (c Container) NotFoundErrorString() string {
	return i18n.ContainerNotFound
}

func (c Container) UpdateErrorString() string {
	return i18n.UpdateContainerError
}

func (c Container) GetUniqueKey() []string {
	return []string{"id"}
}
