package model

import "CBCTF/internel/i18n"

// DockerGroup
// BelongsTo Challenge
// HasMany Docker
type DockerGroup struct {
	ChallengeID     uint            `json:"challenge_id"`
	Challenge       Challenge       `json:"-"`
	Dockers         []Docker        `json:"-"`
	NetworkPolicies NetworkPolicies `gorm:"type:json" json:"network_policies"`
	Basic
}

func (c DockerGroup) GetModelName() string {
	return "DockerGroup"
}

func (c DockerGroup) GetID() uint {
	return c.ID
}

func (c DockerGroup) GetVersion() uint {
	return c.Version
}

func (c DockerGroup) CreateErrorString() string {
	return i18n.CreateDockerGroupError
}

func (c DockerGroup) DeleteErrorString() string {
	return i18n.DeleteDockerGroupError
}

func (c DockerGroup) GetErrorString() string {
	return i18n.GetDockerGroupError
}

func (c DockerGroup) NotFoundErrorString() string {
	return i18n.DockerGroupNotFound
}

func (c DockerGroup) UpdateErrorString() string {
	return i18n.UpdateDockerGroupError
}

func (c DockerGroup) GetUniqueKey() []string {
	return []string{"id"}
}
