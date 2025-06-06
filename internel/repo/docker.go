package repo

import (
	"CBCTF/internel/i18n"
	"CBCTF/internel/log"
	"CBCTF/internel/model"
	"gorm.io/gorm"
)

type DockerRepo struct {
	Basic[model.Docker]
}

type CreateDockerOptions struct {
	DockerGroupID  uint
	Name           string
	Image          string
	PullPolicy     *string
	Hostname       *string
	WorkingDir     *string
	User           *string
	Command        *model.StringList
	Entrypoint     *model.StringList
	CPUCount       *int64
	CPUPercent     *float32
	CPUPeriod      *int64
	CPUQuota       *int64
	CPURTPeriod    *int64
	CPURTRuntime   *int64
	CPUS           *float32
	CPUSet         *string
	CPUShares      *int64
	MemLimit       *int64
	MemReservation *int64
	MemSwapLimit   *int64
	MemSwappiness  *int64
	Expose         *model.StringList
	Environment    *model.StringMap
}

func (c CreateDockerOptions) Convert2Model() model.Model {
	return model.Docker{
		DockerGroupID:  c.DockerGroupID,
		Name:           c.Name,
		Image:          c.Image,
		PullPolicy:     c.PullPolicy,
		Hostname:       c.Hostname,
		WorkingDir:     c.WorkingDir,
		User:           c.User,
		Command:        c.Command,
		Entrypoint:     c.Entrypoint,
		CPUCount:       c.CPUCount,
		CPUPercent:     c.CPUPercent,
		CPUPeriod:      c.CPUPeriod,
		CPUQuota:       c.CPUQuota,
		CPURTPeriod:    c.CPURTPeriod,
		CPURTRuntime:   c.CPURTRuntime,
		CPUS:           c.CPUS,
		CPUSet:         c.CPUSet,
		CPUShares:      c.CPUShares,
		MemLimit:       c.MemLimit,
		MemReservation: c.MemReservation,
		MemSwapLimit:   c.MemSwapLimit,
		MemSwappiness:  c.MemSwappiness,
		Expose:         c.Expose,
		Environment:    c.Environment,
	}
}

type UpdateDockerOptions struct {
	DockerGroupID  *uint
	Name           *string
	Image          *string
	PullPolicy     *string
	Hostname       *string
	WorkingDir     *string
	User           *string
	Command        *model.StringList
	Entrypoint     *model.StringList
	CPUCount       *int64
	CPUPercent     *float32
	CPUPeriod      *int64
	CPUQuota       *int64
	CPURTPeriod    *int64
	CPURTRuntime   *int64
	CPUS           *float32
	CPUSet         *string
	CPUShares      *int64
	MemLimit       *int64
	MemReservation *int64
	MemSwapLimit   *int64
	MemSwappiness  *int64
	Expose         *model.StringList
	Environment    *model.StringMap
}

func (u UpdateDockerOptions) Convert2Map() map[string]any {
	options := make(map[string]any)
	if u.DockerGroupID != nil {
		options["docker_group_id"] = *u.DockerGroupID
	}
	if u.Name != nil {
		options["name"] = *u.Name
	}
	if u.Image != nil {
		options["image"] = *u.Image
	}
	if u.PullPolicy != nil {
		options["pull_policy"] = *u.PullPolicy
	}
	if u.Hostname != nil {
		options["hostname"] = *u.Hostname
	}
	if u.WorkingDir != nil {
		options["working_dir"] = *u.WorkingDir
	}
	if u.User != nil {
		options["user"] = *u.User
	}
	if u.Command != nil {
		options["command"] = *u.Command
	}
	if u.Entrypoint != nil {
		options["entrypoint"] = *u.Entrypoint
	}
	if u.CPUCount != nil {
		options["cpu_count"] = *u.CPUCount
	}
	if u.CPUPercent != nil {
		options["cpu_percent"] = *u.CPUPercent
	}
	if u.CPUPeriod != nil {
		options["cpu_period"] = *u.CPUPeriod
	}
	if u.CPUQuota != nil {
		options["cpu_quota"] = *u.CPUQuota
	}
	if u.CPURTPeriod != nil {
		options["cpu_rt_period"] = *u.CPURTPeriod
	}
	if u.CPURTRuntime != nil {
		options["cpu_rt_runtime"] = *u.CPURTRuntime
	}
	if u.CPUS != nil {
		options["cpus"] = *u.CPUS
	}
	if u.CPUSet != nil {
		options["cpu_set"] = *u.CPUSet
	}
	if u.CPUShares != nil {
		options["cpu_shares"] = *u.CPUShares
	}
	if u.MemLimit != nil {
		options["mem_limit"] = *u.MemLimit
	}
	if u.MemReservation != nil {
		options["mem_reservation"] = *u.MemReservation
	}
	if u.MemSwapLimit != nil {
		options["mem_swap_limit"] = *u.MemSwapLimit
	}
	if u.MemSwappiness != nil {
		options["mem_swappiness"] = *u.MemSwappiness
	}
	if u.Expose != nil {
		options["expose"] = *u.Expose
	}
	if u.Environment != nil {
		options["environment"] = *u.Environment
	}
	return options
}

func InitDockerRepo(tx *gorm.DB) *DockerRepo {
	return &DockerRepo{
		Basic: Basic[model.Docker]{
			DB: tx,
		},
	}
}

func (d *DockerRepo) Delete(idL ...uint) (bool, string) {
	challengeFlagIDL := make([]uint, 0)
	for _, id := range idL {
		docker, ok, msg := d.GetByID(id, "ChallengeFlags")
		if !ok {
			return ok, msg
		}
		for _, challengeFlag := range docker.ChallengeFlags {
			challengeFlagIDL = append(challengeFlagIDL, challengeFlag.ID)
		}
	}
	if ok, msg := InitChallengeFlagRepo(d.DB).Delete(challengeFlagIDL...); !ok {
		return false, msg
	}
	if res := d.DB.Model(&model.Docker{}).Where("id IN ?", idL).Delete(&model.Docker{}); res.Error != nil {
		log.Logger.Warningf("Failed to delete Docker: %v", res.Error)
		return false, model.Docker{}.DeleteErrorString()
	}
	return true, i18n.Success
}
