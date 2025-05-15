package repo

import (
	"CBCTF/internel/model"
	"gorm.io/gorm"
)

type ContainerRepo struct {
	Repo[model.Container]
}

type CreateContainerOptions struct {
	PodID       uint
	Name        string
	Image       string
	Hostname    string
	Flags       model.Strings
	ExposePorts model.Ports
}

func InitContainerRepo(tx *gorm.DB) *ContainerRepo {
	return &ContainerRepo{Repo: Repo[model.Container]{DB: tx, Model: "Container"}}
}
