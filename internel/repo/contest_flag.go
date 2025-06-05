package repo

import (
	"CBCTF/internel/model"
	"gorm.io/gorm"
	"time"
)

type ContestFlagRepo struct {
	Basic[model.ContestFlag]
}

type CreateContestFlagOptions struct {
	ContestID          uint
	ContestChallengeID uint
	Value              string
	Score              float64
	CurrentScore       float64
	Decay              float64
	MinScore           float64
	ScoreType          uint
	Solvers            int64
	Last               time.Time
}

func (c CreateContestFlagOptions) Convert2Model() model.Model {
	return model.ContestFlag{
		ContestID:          c.ContestID,
		ContestChallengeID: c.ContestChallengeID,
		Value:              c.Value,
		Score:              c.Score,
		CurrentScore:       c.CurrentScore,
		Decay:              c.Decay,
		MinScore:           c.MinScore,
		ScoreType:          c.ScoreType,
		Solvers:            c.Solvers,
		Last:               c.Last,
	}
}

type UpdateContestFlagOptions struct {
	Value        *string
	Score        *float64
	CurrentScore *float64
	Decay        *float64
	MinScore     *float64
	ScoreType    *uint
	Solvers      *int64
	Last         *time.Time
}

func (c UpdateContestFlagOptions) Convert2Map() map[string]any {
	options := make(map[string]any)
	if c.Value != nil {
		options["value"] = *c.Value
	}
	if c.Score != nil {
		options["score"] = *c.Score
	}
	if c.CurrentScore != nil {
		options["current_score"] = *c.CurrentScore
	}
	if c.Decay != nil {
		options["decay"] = *c.Decay
	}
	if c.MinScore != nil {
		options["min_score"] = *c.MinScore
	}
	if c.ScoreType != nil {
		options["score_type"] = *c.ScoreType
	}
	if c.Solvers != nil {
		options["solvers"] = *c.Solvers
	}
	if c.Last != nil {
		options["last"] = *c.Last
	}
	return options
}

func InitContestFlagRepo(tx *gorm.DB) *ContestFlagRepo {
	return &ContestFlagRepo{
		Basic: Basic[model.ContestFlag]{
			DB: tx,
		},
	}
}
