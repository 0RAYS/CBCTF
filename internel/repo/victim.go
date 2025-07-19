package repo

import (
	"CBCTF/internel/i18n"
	"CBCTF/internel/log"
	"CBCTF/internel/model"
	"gorm.io/gorm"
	"time"
)

type VictimRepo struct {
	BasicRepo[model.Victim]
}

type CreateVictimOptions struct {
	ContestChallengeID uint
	TeamID             uint
	UserID             uint
	Start              time.Time
	Duration           time.Duration
	VPC                model.VPC
	NetworkPolicies    model.NetworkPolicies
}

func (c CreateVictimOptions) Convert2Model() model.Model {
	return model.Victim{
		ContestChallengeID: c.ContestChallengeID,
		TeamID:             c.TeamID,
		UserID:             c.UserID,
		Start:              c.Start,
		Duration:           c.Duration,
		VPC:                c.VPC,
		NetworkPolicies:    c.NetworkPolicies,
	}
}

type UpdateVictimOptions struct {
	Start    *time.Time
	Duration *time.Duration
}

func (u UpdateVictimOptions) Convert2Map() map[string]any {
	options := make(map[string]any)
	if u.Start != nil {
		options["start"] = *u.Start
	}
	if u.Duration != nil {
		options["duration"] = *u.Duration
	}
	return options
}

func InitVictimRepo(tx *gorm.DB) *VictimRepo {
	return &VictimRepo{
		BasicRepo: BasicRepo[model.Victim]{
			DB: tx,
		},
	}
}

func (v *VictimRepo) HasAliveVictim(teamID uint, contestChallengeID uint) (model.Victim, bool, string) {
	return v.Get(GetOptions{Conditions: map[string]interface{}{"team_id": teamID, "contest_challenge_id": contestChallengeID}})
}

func (v *VictimRepo) Delete(idL ...uint) (bool, string) {
	victimL, _, ok, msg := v.List(-1, -1, GetOptions{
		Conditions: map[string]any{"id": idL},
		Selects:    []string{"id"},
		Preloads: map[string]GetOptions{
			"Pods": {Selects: []string{"id", "victim_id"}},
		},
	})
	if !ok && msg != i18n.VictimNotFound {
		return false, msg
	}
	podIDL := make([]uint, 0)
	for _, victim := range victimL {
		for _, pod := range victim.Pods {
			podIDL = append(podIDL, pod.ID)
		}
	}
	if ok, msg = InitPodRepo(v.DB).Delete(idL...); !ok {
		return false, msg
	}
	if res := v.DB.Model(&model.Victim{}).Where("id IN ?", idL).Delete(&model.Victim{}); res.Error != nil {
		log.Logger.Warningf("Failed to delete Victim: %s", res.Error)
		return false, i18n.DeleteVictimError
	}
	return true, i18n.Success
}
