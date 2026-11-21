package db

import (
	"fmt"
	"time"

	"gorm.io/gorm"

	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
)

type ChallengeFlagRepo struct {
	BaseRepo[model.ChallengeFlag]
}

type UpdateChallengeFlagOptions struct {
	Value        *string
	Binding      *model.FlagBinding
	Generation   *int
	Active       *bool
	Score        *float64
	CurrentScore *float64
	Decay        *float64
	MinScore     *float64
	ScoreType    *model.ScoreType
	Solvers      *int64
	Last         *time.Time
}

func (u UpdateChallengeFlagOptions) Convert2Map() map[string]any {
	options := make(map[string]any)
	if u.Value != nil {
		options["value"] = *u.Value
	}
	if u.Binding != nil {
		options["binding"] = *u.Binding
	}
	if u.Generation != nil {
		options["generation"] = *u.Generation
	}
	if u.Active != nil {
		options["active"] = *u.Active
	}
	if u.Score != nil {
		options["score"] = *u.Score
	}
	if u.CurrentScore != nil {
		options["current_score"] = *u.CurrentScore
	}
	if u.Decay != nil {
		options["decay"] = *u.Decay
	}
	if u.MinScore != nil {
		options["min_score"] = *u.MinScore
	}
	if u.ScoreType != nil {
		options["score_type"] = *u.ScoreType
	}
	if u.Solvers != nil {
		options["solvers"] = *u.Solvers
	}
	if u.Last != nil {
		options["last"] = *u.Last
	}
	return options
}

func InitChallengeFlagRepo(tx *gorm.DB) *ChallengeFlagRepo {
	return &ChallengeFlagRepo{
		DB: tx,
	}
}

func (c *ChallengeFlagRepo) DeleteByChallengeID(challengeIDL ...uint) model.RetVal {
	return c.DeleteByFieldID("challenge_id", challengeIDL...)
}

// ListForContest 返回比赛下所有题目的 flag，题目通过 challenges 关联过滤。
func (c *ChallengeFlagRepo) ListForContest(contestID uint) ([]model.ChallengeFlag, model.RetVal) {
	flags := make([]model.ChallengeFlag, 0)
	res := c.DB.Model(&model.ChallengeFlag{}).
		Joins("INNER JOIN challenges ON challenges.id = challenge_flags.challenge_id AND challenges.deleted_at IS NULL").
		Where("challenges.contest_id = ? AND challenge_flags.deleted_at IS NULL", contestID).
		Preload("Challenge").
		Find(&flags)
	target := fmt.Sprintf("contest_id=%d", contestID)
	if res.Error != nil {
		log.Logger.Warningf("Failed to list ChallengeFlags: %s (target=%s)", res.Error, target)
		return nil, model.FailRetVal(i18n.Model.ChallengeFlag.GetError, "ChallengeFlag", target, res.Error)
	}
	return flags, model.SuccessRetVal()
}

type UserSolvedChallengeFlag struct {
	model.ChallengeFlag
	UserID uint
	TeamID uint
}

type TeamSolvedChallengeFlag struct {
	model.ChallengeFlag
	TeamID uint
}

func (c *ChallengeFlagRepo) GetUserSolvedChallengeFlags(userIDL ...uint) ([]UserSolvedChallengeFlag, model.RetVal) {
	if len(userIDL) == 0 {
		return nil, model.SuccessRetVal()
	}
	var results []UserSolvedChallengeFlag
	res := c.DB.Table("submissions").
		Select("DISTINCT ON (submissions.user_id, challenge_flags.id) submissions.user_id, submissions.team_id, challenge_flags.*").
		Joins("INNER JOIN challenge_flags ON submissions.challenge_flag_id = challenge_flags.id AND challenge_flags.deleted_at IS NULL").
		Joins("INNER JOIN users ON submissions.user_id = users.id AND users.deleted_at IS NULL").
		Joins("INNER JOIN teams ON submissions.team_id = teams.id AND teams.deleted_at IS NULL").
		Where("submissions.user_id IN ? AND submissions.solved = true AND submissions.deleted_at IS NULL", userIDL).
		Order("submissions.user_id, challenge_flags.id, submissions.created_at ASC, submissions.id ASC").
		Scan(&results)
	if res.Error != nil {
		target := fmt.Sprintf("user_id=%v", userIDL)
		log.Logger.Warningf("Failed to get ChallengeFlag: %s (target=%s)", res.Error, target)
		return nil, model.FailRetVal(i18n.Model.ChallengeFlag.GetError, "ChallengeFlag", target, res.Error)
	}
	return results, model.SuccessRetVal()
}

func (c *ChallengeFlagRepo) GetTeamSolvedChallengeFlags(teamIDL ...uint) ([]model.ChallengeFlag, model.RetVal) {
	rows, ret := InitTeamFlagRepo(c.DB).GetSolvedChallengeFlags(teamIDL...)
	if !ret.OK {
		return nil, ret
	}
	results := make([]model.ChallengeFlag, 0, len(rows))
	for _, row := range rows {
		results = append(results, row.ChallengeFlag)
	}
	return results, model.SuccessRetVal()
}

func (c *ChallengeFlagRepo) GetTeamsSolvedChallengeFlags(teamIDL ...uint) ([]TeamSolvedChallengeFlag, model.RetVal) {
	rows, ret := InitTeamFlagRepo(c.DB).GetSolvedChallengeFlags(teamIDL...)
	if !ret.OK {
		return nil, ret
	}
	results := make([]TeamSolvedChallengeFlag, 0, len(rows))
	for _, row := range rows {
		results = append(results, TeamSolvedChallengeFlag{
			TeamID:        row.TeamID,
			ChallengeFlag: row.ChallengeFlag,
		})
	}
	return results, model.SuccessRetVal()
}

// Delete 删除 flag 及其解题记录。
func (c *ChallengeFlagRepo) Delete(idL ...uint) model.RetVal {
	if len(idL) == 0 {
		return model.SuccessRetVal()
	}
	if ret := InitSubmissionRepo(c.DB).DeleteByChallengeFlagID(idL...); !ret.OK {
		return ret
	}
	if ret := InitTeamFlagRepo(c.DB).DeleteByChallengeFlagID(idL...); !ret.OK {
		return ret
	}
	if res := c.DB.Model(&model.ChallengeFlag{}).Where("id IN ?", idL).Delete(&model.ChallengeFlag{}); res.Error != nil {
		target := fmt.Sprintf("id=%v", idL)
		log.Logger.Warningf("Failed to delete ChallengeFlags: %s (target=%s)", res.Error, target)
		return model.FailRetVal(i18n.Model.ChallengeFlag.DeleteError, "ChallengeFlag", target, res.Error)
	}
	return model.SuccessRetVal()
}
