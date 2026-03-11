package db

import (
	"CBCTF/internal/model"
	"time"

	"gorm.io/gorm"
)

type GeneratorRepo struct {
	BaseRepo[model.Generator]
}

type CreateGeneratorOptions struct {
	ChallengeID uint
	ContestID   uint
	Name        string
	Period      time.Duration
}

func (c CreateGeneratorOptions) Convert2Model() model.Model {
	return model.Generator{
		ChallengeID: c.ChallengeID,
		ContestID:   c.ContestID,
		Name:        c.Name,
		Period:      c.Period,
	}
}

type UpdateGeneratorOptions struct {
	Name        *string
	Success     *int64
	SuccessLast *time.Time
	Failure     *int64
	FailureLast *time.Time
	Period      *time.Duration
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
	if u.Period != nil {
		options["period"] = u.Period
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
