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
	ChallengeFlags []ChallengeFlag `json:"-"`
	Name           string          `json:"name"`
	Image          string          `json:"image"`
	PullPolicy     *string         `gorm:"default:null" json:"pull_policy"`
	Hostname       *string         `gorm:"default:null" json:"hostname"`
	WorkingDir     *string         `gorm:"default:null" json:"working_dir"`
	User           *string         `gorm:"default:null" json:"user"`
	Command        *StringList     `gorm:"default:null;type:json" json:"command"`
	Entrypoint     *StringList     `gorm:"default:null;type:json" json:"entrypoint"`
	CPUCount       *int64          `gorm:"default:null" json:"cpu_count"`
	CPUPercent     *float32        `gorm:"default:null" json:"cpu_percent"`
	CPUPeriod      *int64          `gorm:"default:null" json:"cpu_period"`
	CPUQuota       *int64          `gorm:"default:null" json:"cpu_quota"`
	CPURTPeriod    *int64          `gorm:"default:null" json:"cpu_rt_period"`
	CPURTRuntime   *int64          `gorm:"default:null" json:"cpu_rt_runtime"`
	CPUS           *float32        `gorm:"default:null" json:"cpus"`
	CPUSet         *string         `gorm:"default:null" json:"cpu_set"`
	CPUShares      *int64          `gorm:"default:null" json:"cpu_shares"`
	MemLimit       *int64          `gorm:"default:null" json:"mem_limit"`
	MemReservation *int64          `gorm:"default:null" json:"mem_reservation"`
	MemSwapLimit   *int64          `gorm:"default:null" json:"mem_swap_limit"`
	MemSwappiness  *int64          `gorm:"default:null" json:"mem_swappiness"`
	Expose         *StringList     `gorm:"default:null;type:json" json:"expose"`
	Environment    *StringMap      `gorm:"default:null;type:json" json:"environment"`
	Basic
}

func (d Docker) GetModelName() string {
	return "Docker"
}

func (d Docker) GetID() uint {
	return d.ID
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
