package db

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"time"

	"gorm.io/gorm"
)

type ContestFlagRepo struct {
	BaseRepo[model.ContestFlag]
}

type CreateContestFlagOptions struct {
	ContestID          uint
	ContestChallengeID uint
	ChallengeFlagID    uint
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
		ChallengeFlagID:    c.ChallengeFlagID,
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

type DiffUpdateContestFlagOptions struct {
	CurrentScore float64
	Solvers      int64
}

func (d DiffUpdateContestFlagOptions) Convert2Expr() map[string]any {
	options := make(map[string]any)
	if d.CurrentScore != 0 {
		options["current_score"] = gorm.Expr("current_score + ?", d.CurrentScore)
	}
	if d.Solvers != 0 {
		options["solvers"] = gorm.Expr("solvers + ?", d.Solvers)
	}
	return options
}

func InitContestFlagRepo(tx *gorm.DB) *ContestFlagRepo {
	return &ContestFlagRepo{
		BaseRepo: BaseRepo[model.ContestFlag]{
			DB: tx,
		},
	}
}

func (c *ContestFlagRepo) Delete(idL ...uint) (bool, string) {
	contestFlagL, _, ok, msg := c.List(-1, -1, GetOptions{
		Conditions: map[string]any{"id": idL},
		Selects:    []string{"id"},
		Preloads: map[string]GetOptions{
			"Submissions": {Selects: []string{"id", "contest_flag_id"}},
			"TeamFlags":   {Selects: []string{"id", "contest_flag_id"}},
		},
	})
	if !ok {
		if msg != i18n.ContestFlagNotFound {
			return false, msg
		}
		return true, i18n.Success
	}
	submissionIDL, teamFlagIDL := make([]uint, 0), make([]uint, 0)
	for _, contestFlag := range contestFlagL {
		for _, submission := range contestFlag.Submissions {
			submissionIDL = append(submissionIDL, submission.ID)
		}
		for _, teamFlag := range contestFlag.TeamFlags {
			teamFlagIDL = append(teamFlagIDL, teamFlag.ID)
		}
	}
	if ok, msg = InitSubmissionRepo(c.DB).Delete(submissionIDL...); !ok {
		return false, msg
	}
	if ok, msg = InitTeamFlagRepo(c.DB).Delete(teamFlagIDL...); !ok {
		return false, msg
	}
	if res := c.DB.Model(&model.ContestFlag{}).Where("id IN ?", idL).Delete(&model.ContestFlag{}); res.Error != nil {
		log.Logger.Warningf("Failed to delete ContestFlags: %s", res.Error)
		return false, i18n.DeleteContestFlagError
	}
	return true, i18n.Success
}
