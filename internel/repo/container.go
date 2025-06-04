package repo

import (
	"CBCTF/internel/i18n"
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
	Flags       model.StringList
	ExposePorts model.PortList
}

func InitContainerRepo(tx *gorm.DB) *ContainerRepo {
	return &ContainerRepo{
		Repo: Repo[model.Container]{
			DB: tx, Model: "Container",
			CreateError:   i18n.CreateContainerError,
			DeleteError:   i18n.DeleteContainerError,
			GetError:      i18n.GetContainerError,
			NotFoundError: i18n.ContainerNotFound,
		},
	}
}
