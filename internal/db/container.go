package db

import (
	"CBCTF/internal/model"

	"gorm.io/gorm"
)

type ContainerRepo struct {
	BaseRepo[model.Container]
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

func InitContainerRepo(tx *gorm.DB) *ContainerRepo {
	return &ContainerRepo{
		BaseRepo: BaseRepo[model.Container]{
			DB: tx,
		},
	}
}
