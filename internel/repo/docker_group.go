package repo

import (
	"CBCTF/internel/i18n"
	"CBCTF/internel/log"
	"CBCTF/internel/model"
	"gorm.io/gorm"
)

type DockerGroupRepo struct {
	BasicRepo[model.DockerGroup]
}

type CreateDockerGroupOptions struct {
	ChallengeID     uint
	NetworkPolicies model.NetworkPolicies
}

func (c CreateDockerGroupOptions) Convert2Model() model.Model {
	return model.DockerGroup{
		ChallengeID:     c.ChallengeID,
		NetworkPolicies: c.NetworkPolicies,
	}
}

type UpdateDockerGroupOptions struct {
	NetworkPolicies *model.NetworkPolicies
}

func (c UpdateDockerGroupOptions) Convert2Map() map[string]any {
	m := make(map[string]interface{})
	if c.NetworkPolicies != nil {
		m["network_policies"] = c.NetworkPolicies
	}
	return m
}

func InitDockerGroupRepo(tx *gorm.DB) *DockerGroupRepo {
	return &DockerGroupRepo{
		BasicRepo: BasicRepo[model.DockerGroup]{
			DB: tx,
		},
	}
}

func (d *DockerGroupRepo) Delete(idL ...uint) (bool, string) {
	dockerGroupL, _, ok, msg := d.List(-1, -1, GetOptions{
		Conditions: map[string]any{"id": idL},
		Selects:    []string{"id"},
		Preloads: map[string]GetOptions{
			"Dockers": {Selects: []string{"id"}},
		},
	})
	if !ok && msg != i18n.DockerGroupNotFound {
		return false, msg
	}
	dockerIDL := make([]uint, 0)
	for _, dockerGroup := range dockerGroupL {
		for _, docker := range dockerGroup.Dockers {
			dockerIDL = append(dockerIDL, docker.ID)
		}
	}
	if ok, msg = InitDockerRepo(d.DB).Delete(dockerIDL...); !ok {
		return false, msg
	}
	if res := d.DB.Model(&model.DockerGroup{}).Where("id IN ?", idL).Delete(&model.DockerGroup{}); res.Error != nil {
		log.Logger.Warningf("Failed to delete DockerGroup: %v", res.Error)
		return false, i18n.DeleteDockerGroupError
	}
	return true, i18n.Success
}
