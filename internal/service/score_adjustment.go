package service

import (
	"strings"

	"gorm.io/gorm"

	"CBCTF/internal/db"
	"CBCTF/internal/model"
)

// AdjustTeamScore 记录一次管理员分数变动，并按解题分数加变动总和刷新队伍分数与排名。
func AdjustTeamScore(tx *gorm.DB, contest model.Contest, team model.Team, score float64, reason string, creatorID uint) (model.ScoreAdjustment, model.RetVal) {
	var adjustment model.ScoreAdjustment
	ret := db.WithTransactionDB(tx, func(work *gorm.DB) model.RetVal {
		created, ret := db.InitScoreAdjustmentRepo(work).Create(model.ScoreAdjustment{
			ContestID: contest.ID,
			TeamID:    team.ID,
			Score:     score,
			Reason:    strings.TrimSpace(reason),
			CreatorID: creatorID,
		})
		if !ret.OK {
			return ret
		}
		adjustment = created
		if _, _, ret := UpdateTeamRanking(work, contest, -1, -1); !ret.OK {
			return ret
		}
		return model.SuccessRetVal()
	})
	if !ret.OK {
		return model.ScoreAdjustment{}, ret
	}
	return adjustment, model.SuccessRetVal()
}

// ListTeamScoreAdjustments 返回队伍的分数变动记录。
func ListTeamScoreAdjustments(tx *gorm.DB, contest model.Contest, team model.Team) ([]model.ScoreAdjustment, model.RetVal) {
	return db.InitScoreAdjustmentRepo(tx).ListByTeamID(contest.ID, team.ID)
}
