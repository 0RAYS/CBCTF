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
	UserID             sql.Null[uint]
	Start              time.Time
	Duration           time.Duration
	VPC                model.VPC
	NetworkPolicies    model.NetworkPolicies
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
		VPC:                c.VPC,
		NetworkPolicies:    c.NetworkPolicies,
	}
}

type UpdateVictimOptions struct {
	Start            *time.Time
	Duration         *time.Duration
	VPC              *model.VPC
	Endpoints        *model.Endpoints
	ExposedEndpoints *model.Endpoints
}

func (u UpdateVictimOptions) Convert2Map() map[string]any {
	options := make(map[string]any)
	if u.Start != nil {
		options["start"] = *u.Start
	}
	if u.Duration != nil {
		options["duration"] = *u.Duration
	}
	if u.VPC != nil {
		options["vpc"] = *u.VPC
	}
	if u.Endpoints != nil {
		options["endpoints"] = *u.Endpoints
	}
	if u.ExposedEndpoints != nil {
		options["exposed_endpoints"] = *u.ExposedEndpoints
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
	victimL, _, ret := v.List(-1, -1, GetOptions{
		Conditions: map[string]any{"id": idL},
		Selects:    []string{"id"},
		Preloads:   map[string]GetOptions{"Pods": {Selects: []string{"id", "victim_id"}}},
	})
	if !ret.OK {
		if ret.Msg != i18n.Model.NotFound {
			return ret
		}
		return model.SuccessRetVal()
	}
	podIDL := make([]uint, 0)
	for _, victim := range victimL {
		for _, pod := range victim.Pods {
			podIDL = append(podIDL, pod.ID)
		}
	}
	if ret = InitPodRepo(v.DB).Delete(podIDL...); !ret.OK {
		return ret
	}
	if res := v.DB.Model(&model.Victim{}).Where("id IN ?", idL).Delete(&model.Victim{}); res.Error != nil {
		log.Logger.Warningf("Failed to delete Victim: %s", res.Error)
		return model.RetVal{Msg: i18n.Model.DeleteError, Attr: map[string]interface{}{"Model": model.Victim{}.ModelName(), "Error": res.Error.Error()}}
	}
	return model.SuccessRetVal()
}
