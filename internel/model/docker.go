package model

import "CBCTF/internel/i18n"

const (
	EnvFlagPrefix      = "FLAG"
	VolumeFlagPrefix   = "FLAG"
	VolumeFlagLabelKey = "value"
)

// Docker
// BelongsTo DockerGroup
// HasMany ChallengeFlag
type Docker struct {
	DockerGroupID  uint            `json:"docker_group_id"`
	DockerGroup    DockerGroup     `json:"-"`
	ChallengeFlags []ChallengeFlag `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	Name           string          `json:"name"`
	Image          string          `json:"image"`
	WorkingDir     *string         `gorm:"default:null" json:"working_dir"`
	Command        StringList      `gorm:"default:null;type:json" json:"command"`
	Expose         StringList      `gorm:"default:null;type:json" json:"expose"`
	Environment    StringMap       `gorm:"default:null;type:json" json:"environment"`
	BasicModel
}

func (d Docker) GetModelName() string {
	return "Docker"
}

func (d Docker) GetVersion() uint {
	return d.Version
}

func (d Docker) CreateErrorString() string {
	return i18n.CreateDockerError
}

func (d Docker) DeleteErrorString() string {
	return i18n.DeleteDockerError
}

func (d Docker) GetErrorString() string {
	return i18n.GetDockerError
}

func (d Docker) NotFoundErrorString() string {
	return i18n.DockerNotFound
}

func (d Docker) UpdateErrorString() string {
	return i18n.UpdateDockerError
}

func (d Docker) GetUniqueKey() []string {
	return []string{"id"}
}
