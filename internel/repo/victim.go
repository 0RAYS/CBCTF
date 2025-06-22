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
	IPBlock            string
	Start              time.Time
	Duration           time.Duration
	HostAlias          model.StringMap
}

func (c CreateVictimOptions) Convert2Model() model.Model {
	return model.Victim{
		ContestChallengeID: c.ContestChallengeID,
		TeamID:             c.TeamID,
		UserID:             c.UserID,
		IPBlock:            c.IPBlock,
		Start:              c.Start,
		Duration:           c.Duration,
		HostAlias:          c.HostAlias,
	}
}

type UpdateVictimOptions struct {
	IPBlock  *string
	Start    *time.Time
	Duration *time.Duration
}

func (u UpdateVictimOptions) Convert2Map() map[string]any {
	options := make(map[string]any)
	if u.IPBlock != nil {
		options["ip_block"] = *u.IPBlock
	}
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
	podIDL := make([]uint, 0)
	for _, id := range idL {
		victim, ok, msg := v.GetByID(id, GetOptions{
			Selects: []string{"id"},
			Preloads: map[string]GetOptions{
				"Pods": {Selects: []string{"id"}},
			},
		})
		if !ok && msg != i18n.VictimNotFound {
			return false, msg
		}
		for _, pod := range victim.Pods {
			podIDL = append(podIDL, pod.ID)
		}
	}
	if ok, msg := InitPodRepo(v.DB).Delete(idL...); !ok {
		return false, msg
	}
	if res := v.DB.Model(&model.Victim{}).Where("id IN ?", idL).Delete(&model.Victim{}); res.Error != nil {
		log.Logger.Warningf("Failed to delete Victim: %s", res.Error)
		return false, i18n.DeleteVictimError
	}
	return true, i18n.Success
}
