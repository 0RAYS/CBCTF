package repo

import (
	"CBCTF/internal/model"
	"gorm.io/gorm"
)

type ContainerRepo struct {
	BasicRepo[model.Container]
}

type CreateContainerOptions struct {
	PodID       uint
	Name        string
	Image       string
	CPU         float32
	Memory      int64
	WorkingDir  string
	Command     model.StringList
	Environment model.StringMap
	EnvFlags    model.StringMap
	VolumeFlags model.StringMap
	Exposes     model.Exposes
}

func (c CreateContainerOptions) Convert2Model() model.Model {
	return model.Container{
		PodID:       c.PodID,
		Name:        c.Name,
		Image:       c.Image,
		CPU:         c.CPU,
		Memory:      c.Memory,
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
