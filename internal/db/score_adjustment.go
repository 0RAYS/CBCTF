package db

import (
	"fmt"

	"gorm.io/gorm"

	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
)

type ScoreAdjustmentRepo struct {
	BaseRepo[model.ScoreAdjustment]
}

func InitScoreAdjustmentRepo(tx *gorm.DB) *ScoreAdjustmentRepo {
	return &ScoreAdjustmentRepo{
		DB: tx,
	}
}

// SumByTeamID 汇总队伍在比赛中的分数变动总和，没有变动记录的队伍不在结果中。
func (r *ScoreAdjustmentRepo) SumByTeamID(contestID uint, teamIDL ...uint) (map[uint]float64, model.RetVal) {
	sums := make(map[uint]float64, len(teamIDL))
	if len(teamIDL) == 0 {
		return sums, model.SuccessRetVal()
	}
	type sumRow struct {
		TeamID uint
		Total  float64
	}
	var rows []sumRow
	res := r.DB.Model(&model.ScoreAdjustment{}).
		Select("team_id, COALESCE(SUM(score), 0) AS total").
		Where("contest_id = ? AND team_id IN ?", contestID, teamIDL).
		Group("team_id").
		Scan(&rows)
	if res.Error != nil {
		log.Logger.Warningf("Failed to sum score adjustments: %s (contest_id=%d)", res.Error, contestID)
		return nil, model.FailRetVal(i18n.Model.GetError, "ScoreAdjustment", fmt.Sprintf("contest_id=%d", contestID), res.Error)
	}
	for _, row := range rows {
		sums[row.TeamID] = row.Total
	}
	return sums, model.SuccessRetVal()
}

// ListByTeamID 返回队伍在比赛中的分数变动记录，最新的排在前面。
func (r *ScoreAdjustmentRepo) ListByTeamID(contestID, teamID uint) ([]model.ScoreAdjustment, model.RetVal) {
	adjustments, _, ret := r.List(-1, -1, GetOptions{
		Conditions: map[string]any{"contest_id": contestID, "team_id": teamID},
		Sort:       []string{"id DESC"},
	})
	return adjustments, ret
}
