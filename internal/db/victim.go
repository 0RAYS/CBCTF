package db

import (
	"fmt"
	"time"

	"gorm.io/gorm"

	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
)

type VictimRepo struct {
	BaseRepo[model.Victim]
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
		DB: tx,
	}
}

// Create skips generic uniqueness preflight queries. Victims use a database
// generated primary key and have no natural unique key to validate here.
func (v *VictimRepo) Create(victim model.Victim) (model.Victim, model.RetVal) {
	if res := v.DB.Model(&model.Victim{}).Create(&victim); res.Error != nil {
		log.Logger.Warningf("Failed to create Victim: %s", res.Error)
		return model.Victim{}, model.FailRetVal(i18n.Model.CreateError, model.Name(model.Victim{}), "", res.Error)
	}
	return victim, model.SuccessRetVal()
}

// UpdateIfStatus updates startup-owned fields directly only if the victim is
// still in the expected state. This avoids the generic read-before-write while
// preventing startup from overwriting a concurrent stop transition.
func (v *VictimRepo) UpdateIfStatus(id uint, expectedStatus string, options UpdateVictimOptions) model.RetVal {
	data := options.Convert2Map()
	if len(data) == 0 {
		return model.SuccessRetVal()
	}
	res := v.DB.Model(&model.Victim{}).Where("id = ? AND status = ?", id, expectedStatus).Updates(data)
	target := fmt.Sprintf("id=%d,status=%s", id, expectedStatus)
	if res.Error != nil {
		log.Logger.Warningf("Failed to update Victim: %s (target=%s)", res.Error, target)
		return model.FailRetVal(i18n.Model.UpdateError, model.Name(model.Victim{}), target, res.Error)
	}
	if res.RowsAffected == 0 {
		return model.FailRetVal(i18n.Model.Victim.NotStartable, model.Name(model.Victim{}), target, nil).
			WithAttrs(map[string]any{"ExpectedStatus": expectedStatus})
	}
	return model.SuccessRetVal()
}

func (v *VictimRepo) HasAliveVictim(teamID, challengeID uint) (model.Victim, model.RetVal) {
	options := GetOptions{Conditions: map[string]any{"team_id": nil, "challenge_id": challengeID}}
	if teamID > 0 {
		options.Conditions["team_id"] = teamID
	}
	return v.Get(options)
}

// HasAliveVictimByContestChallenge 按比赛题目查找存活的靶机。
// 同一题库题目允许在同场比赛多次加入，靶机唯一性以比赛题目为准。
func (v *VictimRepo) HasAliveVictimByContestChallenge(teamID, contestChallengeID uint) (model.Victim, model.RetVal) {
	options := GetOptions{Conditions: map[string]any{"team_id": nil, "contest_challenge_id": contestChallengeID}}
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
		target := fmt.Sprintf("id=%v", idL)
		log.Logger.Warningf("Failed to delete Victim: %s (target=%s)", res.Error, target)
		return model.FailRetVal(i18n.Model.Victim.DeleteError, model.Name(model.Victim{}), target, res.Error)
	}
	return model.SuccessRetVal()
}

func (v *VictimRepo) DeleteByChallengeID(challengeIDL ...uint) model.RetVal {
	return v.DeleteByFieldID("challenge_id", challengeIDL...)
}

func (v *VictimRepo) DeleteByContestID(contestIDL ...uint) model.RetVal {
	return v.DeleteByFieldID("contest_id", contestIDL...)
}
