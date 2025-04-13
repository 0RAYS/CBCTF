package repo

import "CBCTF/internel/model"

type ContainerRepo struct {
	Repo[model.Container]
}

type CreateContainerOptions struct {
	PodID    uint
	Name     string
	Image    string
	Hostname string
	Flags    model.Strings
}

func InitContainerRepo() *ContainerRepo {
	return &ContainerRepo{Repo: Repo[model.Container]{DB: DB, Model: "Container"}}
}
