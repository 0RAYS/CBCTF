package db

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"database/sql"
	"time"

	"gorm.io/gorm"
)

type GeneratorRepo struct {
	BaseRepo[model.Generator]
}

type CreateGeneratorOptions struct {
	ChallengeID   uint
	ChallengeName string
	ContestID     sql.Null[uint]
	Name          string
}

func (c CreateGeneratorOptions) Convert2Model() model.Model {
	return model.Generator{
		ChallengeID:   c.ChallengeID,
		ChallengeName: c.ChallengeName,
		ContestID:     c.ContestID,
		Name:          c.Name,
		Status:        model.WaitingGeneratorStatus,
	}
}

type UpdateGeneratorOptions struct {
	Name        *string
	Success     *int64
	SuccessLast *time.Time
	Failure     *int64
	FailureLast *time.Time
	Status      *string
}

func (u UpdateGeneratorOptions) Convert2Map() map[string]any {
	options := make(map[string]any)
	if u.Name != nil {
		options["name"] = u.Name
	}
	if u.Success != nil {
		options["success"] = u.Success
	}
	if u.SuccessLast != nil {
		options["success_last"] = u.SuccessLast
	}
	if u.Failure != nil {
		options["failure"] = u.Failure
	}
	if u.FailureLast != nil {
		options["failure_last"] = u.FailureLast
	}
	if u.Status != nil {
		options["status"] = u.Status
	}
	return options
}

type DiffUpdateGeneratorOptions struct {
	Success int64
	Failure int64
}

func (d DiffUpdateGeneratorOptions) Convert2Expr() map[string]any {
	options := make(map[string]any)
	if d.Success != 0 {
		options["success"] = gorm.Expr("success + ?", d.Success)
	}
	if d.Failure != 0 {
		options["failure"] = gorm.Expr("failure + ?", d.Failure)
	}
	return options
}

func InitGeneratorRepo(tx *gorm.DB) *GeneratorRepo {
	return &GeneratorRepo{
		BaseRepo: BaseRepo[model.Generator]{
			DB: tx,
		},
	}
}

func (g *GeneratorRepo) UpdateStatus(id uint, success bool, last time.Time) model.RetVal {
	var diffOptions DiffUpdateGeneratorOptions
	var options UpdateGeneratorOptions
	if success {
		diffOptions = DiffUpdateGeneratorOptions{
			Success: 1,
		}
		options = UpdateGeneratorOptions{
			SuccessLast: &last,
		}
	} else {
		diffOptions = DiffUpdateGeneratorOptions{
			Failure: 1,
		}
		options = UpdateGeneratorOptions{
			FailureLast: &last,
		}
	}
	if ret := g.DiffUpdate(id, diffOptions); !ret.OK {
		return ret
	}
	return g.Update(id, options)
}

func (g *GeneratorRepo) Delete(idL ...uint) model.RetVal {
	for _, id := range idL {
		if ret := g.Update(id, UpdateGeneratorOptions{Status: new(model.StoppedGeneratorStatus)}); !ret.OK {
			return ret
		}
	}
	if res := g.DB.Model(&model.Generator{}).Where("id IN ?", idL).Delete(&model.Generator{}); res.Error != nil {
		log.Logger.Warningf("Failed to delete Generator: %s", res.Error)
		return model.RetVal{Msg: i18n.Model.Generator.DeleteError, Attr: map[string]any{"Error": res.Error.Error()}}
	}
	return model.SuccessRetVal()
}

func (g *GeneratorRepo) DeleteByChallengeID(challengeIDL ...uint) model.RetVal {
	if len(challengeIDL) == 0 {
		return model.SuccessRetVal()
	}
	var generatorIDL []uint
	if res := g.DB.Model(&model.Generator{}).Where("challenge_id IN ?", challengeIDL).Pluck("id", &generatorIDL); res.Error != nil {
		log.Logger.Warningf("Failed to get Generators by challenge IDs %v: %s", challengeIDL, res.Error)
		return model.RetVal{Msg: i18n.Model.Generator.DeleteError, Attr: map[string]any{"Error": res.Error.Error()}}
	}
	if len(generatorIDL) == 0 {
		return model.SuccessRetVal()
	}
	return g.Delete(generatorIDL...)
}

func (g *GeneratorRepo) DeleteByContestID(contestIDL ...uint) model.RetVal {
	if len(contestIDL) == 0 {
		return model.SuccessRetVal()
	}
	var generatorIDL []uint
	if res := g.DB.Model(&model.Generator{}).Where("contest_id IN ?", contestIDL).Pluck("id", &generatorIDL); res.Error != nil {
		log.Logger.Warningf("Failed to get Generators by contest IDs %v: %s", contestIDL, res.Error)
		return model.RetVal{Msg: i18n.Model.Generator.DeleteError, Attr: map[string]any{"Error": res.Error.Error()}}
	}
	if len(generatorIDL) == 0 {
		return model.SuccessRetVal()
	}
	return g.Delete(generatorIDL...)
}
