package db

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"

	"gorm.io/gorm"
)

type DockerRepo struct {
	BaseRepo[model.Docker]
}

type CreateDockerOptions struct {
	ChallengeID uint
	Name        string
	Image       string
	CPU         float32
	Memory      int64
	WorkingDir  string
	Command     model.StringList
	Exposes     model.Exposes
	Environment model.StringMap
	Networks    model.Networks
}

func (c CreateDockerOptions) Convert2Model() model.Model {
	return model.Docker{
		ChallengeID: c.ChallengeID,
		Name:        c.Name,
		Image:       c.Image,
		CPU:         c.CPU,
		Memory:      c.Memory,
		WorkingDir:  c.WorkingDir,
		Command:     c.Command,
		Exposes:     c.Exposes,
		Environment: c.Environment,
		Networks:    c.Networks,
	}
}

func InitDockerRepo(tx *gorm.DB) *DockerRepo {
	return &DockerRepo{
		BaseRepo: BaseRepo[model.Docker]{
			DB: tx,
		},
	}
}

func (d *DockerRepo) Delete(idL ...uint) model.RetVal {
	dockerL, _, ret := d.List(-1, -1, GetOptions{
		Conditions: map[string]any{"id": idL},
		Preloads: map[string]GetOptions{
			"ChallengeFlags": {},
		},
	})
	if !ret.OK {
		if ret.Msg != i18n.Model.NotFound {
			return ret
		}
		return model.SuccessRetVal()
	}
	challengeFlagIDL := make([]uint, 0)
	for _, docker := range dockerL {
		for _, challengeFlag := range docker.ChallengeFlags {
			challengeFlagIDL = append(challengeFlagIDL, challengeFlag.ID)
		}
	}
	if ret = InitChallengeFlagRepo(d.DB).Delete(challengeFlagIDL...); !ret.OK {
		return ret
	}
	if res := d.DB.Model(&model.Docker{}).Where("id IN ?", idL).Delete(&model.Docker{}); res.Error != nil {
		log.Logger.Warningf("Failed to delete Docker: %s", res.Error)
		return model.RetVal{Msg: i18n.Model.Docker.DeleteError, Attr: map[string]any{"Error": res.Error.Error()}}
	}
	return model.SuccessRetVal()
}
