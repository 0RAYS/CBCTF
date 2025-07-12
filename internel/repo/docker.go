package repo

import (
	"CBCTF/internel/i18n"
	"CBCTF/internel/log"
	"CBCTF/internel/model"
	"gorm.io/gorm"
)

type DockerRepo struct {
	BasicRepo[model.Docker]
}

type CreateDockerOptions struct {
	ChallengeID uint
	Name        string
	Image       string
	CPU         float32
	Memory      int64
	WorkingDir  string
	Command     model.StringList
	Expose      model.StringList
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
		Expose:      c.Expose,
		Environment: c.Environment,
		Networks:    c.Networks,
	}
}

type UpdateDockerOptions struct {
}

func (u UpdateDockerOptions) Convert2Map() map[string]any {
	options := make(map[string]any)
	return options
}

func InitDockerRepo(tx *gorm.DB) *DockerRepo {
	return &DockerRepo{
		BasicRepo: BasicRepo[model.Docker]{
			DB: tx,
		},
	}
}

func (d *DockerRepo) Delete(idL ...uint) (bool, string) {
	dockerL, _, ok, msg := d.List(-1, -1, GetOptions{
		Conditions: map[string]any{"id": idL},
		Selects:    []string{"id"},
		Preloads: map[string]GetOptions{
			"ChallengeFlags": {Selects: []string{"id", "docker_id"}},
		},
	})
	if !ok && msg != i18n.DockerNotFound {
		return ok, msg
	}
	challengeFlagIDL := make([]uint, 0)
	for _, docker := range dockerL {
		for _, challengeFlag := range docker.ChallengeFlags {
			challengeFlagIDL = append(challengeFlagIDL, challengeFlag.ID)
		}
	}
	if ok, msg = InitChallengeFlagRepo(d.DB).Delete(challengeFlagIDL...); !ok {
		return false, msg
	}
	if res := d.DB.Model(&model.Docker{}).Where("id IN ?", idL).Delete(&model.Docker{}); res.Error != nil {
		log.Logger.Warningf("Failed to delete Docker: %v", res.Error)
		return false, i18n.DeleteDockerError
	}
	return true, i18n.Success
}
