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
	PullPolicy     *string         `json:"pull_policy"`
	Hostname       *string         `json:"hostname"`
	WorkingDir     *string         `json:"working_dir"`
	User           *string         `json:"user"`
	Command        *StringList     `json:"command"`
	Entrypoint     *StringList     `json:"entrypoint"`
	CPUCount       *int64          `json:"cpu_count"`
	CPUPercent     *float32        `json:"cpu_percent"`
	CPUPeriod      *int64          `json:"cpu_period"`
	CPUQuota       *int64          `json:"cpu_quota"`
	CPURTPeriod    *int64          `json:"cpu_rt_period"`
	CPURTRuntime   *int64          `json:"cpu_rt_runtime"`
	CPUS           *float32        `json:"cpus"`
	CPUSet         *string         `json:"cpu_set"`
	CPUShares      *int64          `json:"cpu_shares"`
	MemLimit       *int64          `json:"mem_limit"`
	MemReservation *int64          `json:"mem_reservation"`
	MemSwapLimit   *int64          `json:"mem_swap_limit"`
	MemSwappiness  *int64          `json:"mem_swappiness"`
	Expose         *StringList     `json:"expose"`
	Environment    *StringMap      `json:"environment"`
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
