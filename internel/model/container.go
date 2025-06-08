package model

import "CBCTF/internel/i18n"

type Container struct {
	PodID       uint        `json:"pod_id"`
	Pod         Pod         `json:"-"`
	Name        string      `json:"name"`
	Image       string      `json:"image"`
	Hostname    string      `json:"hostname"`
	PullPolicy  *string     `gorm:"default:null" json:"pull_policy"`
	WorkingDir  *string     `gorm:"default:null" json:"working_dir"`
	Command     *StringList `gorm:"default:null;type:json" json:"command"`
	Environment *StringMap  `gorm:"default:null;type:json" json:"environment"`
	EnvFlags    StringList  `gorm:"type:json" json:"env_flags"`
	VolumeFlags StringMap   `gorm:"type:json" json:"volume_flags"`
	Exposes     StringList  `gorm:"type:json" json:"exposes"`
	Basic
}

func (c Container) GetModelName() string {
	return "Container"
}

func (c Container) GetID() uint {
	return c.ID
}

func (c Container) GetVersion() uint {
	return c.Version
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
