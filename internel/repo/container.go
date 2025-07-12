package repo

import (
	"CBCTF/internel/model"
	"gorm.io/gorm"
)

type ContainerRepo struct {
	BasicRepo[model.Container]
}

type CreateContainerOptions struct {
	PodID       uint
	Name        string
	Image       string
	Hostname    string
	WorkingDir  string
	Command     model.StringList
	Environment model.StringMap
	EnvFlags    model.StringList
	VolumeFlags model.StringMap
	Exposes     model.StringList
}

func (c CreateContainerOptions) Convert2Model() model.Model {
	return model.Container{
		PodID:       c.PodID,
		Name:        c.Name,
		Image:       c.Image,
		Hostname:    c.Hostname,
		WorkingDir:  c.WorkingDir,
		Command:     c.Command,
		Environment: c.Environment,
		EnvFlags:    c.EnvFlags,
		VolumeFlags: c.VolumeFlags,
		Exposes:     c.Exposes,
	}
}

type UpdateContainerOptions struct {
}

func (u UpdateContainerOptions) Convert2Map() map[string]any {
	options := make(map[string]any)
	return options
}

func InitContainerRepo(tx *gorm.DB) *ContainerRepo {
	return &ContainerRepo{
		BasicRepo: BasicRepo[model.Container]{
			DB: tx,
		},
	}
}
