package db

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"database/sql"
	"time"

	"gorm.io/gorm"
)

type VictimRepo struct {
	BaseRepo[model.Victim]
}

type CreateVictimOptions struct {
	ChallengeID        uint
	ContestID          sql.Null[uint]
	ContestChallengeID sql.Null[uint]
	TeamID             sql.Null[uint]
	UserID             uint
	Start              time.Time
	Duration           time.Duration
	Spec               model.VictimSpec
}

func (c CreateVictimOptions) Convert2Model() model.Model {
	return model.Victim{
		ChallengeID:        c.ChallengeID,
		ContestID:          c.ContestID,
		ContestChallengeID: c.ContestChallengeID,
		TeamID:             c.TeamID,
		UserID:             c.UserID,
		Start:              c.Start,
		Duration:           c.Duration,
		Spec:               c.Spec,
		Status:             model.WaitingVictimStatus,
	}
}

type UpdateVictimOptions struct {
	Start            *time.Time
	Duration         *time.Duration
	Spec             *model.VictimSpec
	Resources        *model.VictimResources
	Endpoints        *model.Endpoints
	ExposedEndpoints *model.Endpoints
	Status           *string
}

func (u UpdateVictimOptions) Convert2Map() map[string]any {
	options := make(map[string]any)
	if u.Start != nil {
		options["start"] = *u.Start
	}
	if u.Duration != nil {
		options["duration"] = *u.Duration
	}
	if u.Spec != nil {
		options["spec"] = *u.Spec
	}
	if u.Resources != nil {
		options["resources"] = *u.Resources
	}
	if u.Endpoints != nil {
		options["endpoints"] = *u.Endpoints
	}
	if u.ExposedEndpoints != nil {
		options["exposed_endpoints"] = *u.ExposedEndpoints
	}
	if u.Status != nil {
		options["status"] = *u.Status
	}
	return options
}

func InitVictimRepo(tx *gorm.DB) *VictimRepo {
	return &VictimRepo{
		BaseRepo: BaseRepo[model.Victim]{
			DB: tx,
		},
	}
}

func (v *VictimRepo) HasAliveVictim(teamID, challengeID uint) (model.Victim, model.RetVal) {
	options := GetOptions{Conditions: map[string]any{"team_id": nil, "challenge_id": challengeID}}
	if teamID > 0 {
		options.Conditions["team_id"] = teamID
	}
	return v.Get(options)
}

func (v *VictimRepo) Delete(idL ...uint) model.RetVal {
	if len(idL) == 0 {
		return model.SuccessRetVal()
	}
	victimL, ret := v.FindAll(GetOptions{
		Conditions: map[string]any{"id": idL},
	})
	if !ret.OK {
		if ret.Msg != i18n.Model.NotFound && ret.Msg != i18n.Model.Victim.NotFound {
			return ret
		}
		return model.SuccessRetVal()
	}
	if ret = InitPodRepo(v.DB).DeleteByVictimID(idL...); !ret.OK {
		return ret
	}
	for _, victim := range victimL {
		if ret = v.Update(victim.ID, UpdateVictimOptions{Status: new(model.StoppedVictimStatus)}); !ret.OK {
			return ret
		}
	}
	if res := v.DB.Model(&model.Victim{}).Where("id IN ?", idL).Delete(&model.Victim{}); res.Error != nil {
		log.Logger.Warningf("Failed to delete Victim: %s", res.Error)
		return model.RetVal{Msg: i18n.Model.Victim.DeleteError, Attr: map[string]any{"Error": res.Error.Error()}}
	}
	return model.SuccessRetVal()
}

func (v *VictimRepo) DeleteByChallengeID(challengeIDL ...uint) model.RetVal {
	return v.DeleteByFieldID("challenge_id", challengeIDL...)
}

func (v *VictimRepo) DeleteByContestID(contestIDL ...uint) model.RetVal {
	return v.DeleteByFieldID("contest_id", contestIDL...)
}
