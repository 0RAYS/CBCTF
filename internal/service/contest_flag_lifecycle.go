package service

import (
	"fmt"

	"gorm.io/gorm"

	"CBCTF/internal/db"
	"CBCTF/internal/dto"
	"CBCTF/internal/model"
)

// UpdateContestFlag 修改比赛题目的 flag，默认保留已解出队伍的解题记录。
// 改 flag 模板时代次 +1，未解出队伍的 team_flags 按新模板重新渲染。
func UpdateContestFlag(tx *gorm.DB, contest model.Contest, contestFlag model.ContestFlag, form dto.UpdateContestFlagForm) model.RetVal {
	currentScore := contestFlag.CurrentScore
	if form.Score != nil && *form.Score < currentScore {
		currentScore = *form.Score
	}
	options := db.UpdateContestFlagOptions{
		Score:        form.Score,
		CurrentScore: &currentScore,
		Decay:        form.Decay,
		MinScore:     form.MinScore,
		ScoreType:    form.ScoreType,
	}
	if form.Value == nil || *form.Value == contestFlag.Value {
		return db.InitContestFlagRepo(tx).Update(contestFlag.ID, options)
	}
	generation := contestFlag.Generation + 1
	options.Value = form.Value
	options.Generation = &generation
	if ret := db.InitContestFlagRepo(tx).Update(contestFlag.ID, options); !ret.OK {
		return ret
	}
	return RerenderTeamFlagValues(tx, contest, contestFlag.ID, *form.Value)
}

// RerenderTeamFlagValues 按新的 flag 模板重新渲染未解出队伍手里的答案，已解出的队伍保留原值。
func RerenderTeamFlagValues(tx *gorm.DB, contest model.Contest, contestFlagID uint, template string) model.RetVal {
	teamFlags, _, ret := db.InitTeamFlagRepo(tx).List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"contest_flag_id": contestFlagID},
	})
	if !ret.OK {
		return ret
	}
	repo := db.InitTeamFlagRepo(tx)
	for _, teamFlag := range teamFlags {
		if teamFlag.Solved {
			continue
		}
		value := renderChallengeFlagValue(template)
		if prefix := contest.Prefix; prefix != "" {
			value = fmt.Sprintf("%s{%s}", prefix, value)
		}
		if teamFlag.Value == value {
			continue
		}
		if ret := repo.Update(teamFlag.ID, db.UpdateTeamFlagRepo{Value: &value}); !ret.OK {
			return ret
		}
	}
	return model.SuccessRetVal()
}

// ResetContestChallengeSolves 清空某比赛题目的解题记录并重算分数与排名。
func ResetContestChallengeSolves(tx *gorm.DB, contest model.Contest, contestChallenge model.ContestChallenge) model.RetVal {
	var clearedTeams, clearedSubmissions int
	ret := db.WithTransactionDB(tx, func(tx2 *gorm.DB) model.RetVal {
		contestFlags, _, ret := db.InitContestFlagRepo(tx2).List(-1, -1, db.GetOptions{
			Conditions: map[string]any{"contest_challenge_id": contestChallenge.ID},
		})
		if !ret.OK {
			return ret
		}
		contestFlagIDL := make([]uint, 0, len(contestFlags))
		for _, contestFlag := range contestFlags {
			contestFlagIDL = append(contestFlagIDL, contestFlag.ID)
		}
		if len(contestFlagIDL) > 0 {
			teamFlags, _, ret := db.InitTeamFlagRepo(tx2).List(-1, -1, db.GetOptions{
				Conditions: map[string]any{"contest_flag_id": contestFlagIDL},
			})
			if !ret.OK {
				return ret
			}
			teams := make(map[uint]struct{}, len(teamFlags))
			for _, teamFlag := range teamFlags {
				teams[teamFlag.TeamID] = struct{}{}
			}
			clearedTeams = len(teams)
			if ret := db.InitTeamFlagRepo(tx2).DeleteByContestFlagID(contestFlagIDL...); !ret.OK {
				return ret
			}
		}
		solvedSubmissions, _, ret := db.InitSubmissionRepo(tx2).List(-1, -1, db.GetOptions{
			Conditions: map[string]any{
				"contest_challenge_id": contestChallenge.ID,
				"solved":               true,
			},
		})
		if !ret.OK {
			return ret
		}
		submissionIDL := make([]uint, 0, len(solvedSubmissions))
		for _, submission := range solvedSubmissions {
			submissionIDL = append(submissionIDL, submission.ID)
		}
		clearedSubmissions = len(submissionIDL)
		if len(submissionIDL) > 0 {
			if ret := db.InitSubmissionRepo(tx2).Delete(submissionIDL...); !ret.OK {
				return ret
			}
		}
		// 重置 flag 的计分状态
		flagRepo := db.InitContestFlagRepo(tx2)
		for _, contestFlag := range contestFlags {
			score := contestFlag.Score
			if ret := flagRepo.Update(contestFlag.ID, db.UpdateContestFlagOptions{
				Solvers:      new(int64(0)),
				CurrentScore: &score,
			}); !ret.OK {
				return ret
			}
		}
		return model.SuccessRetVal()
	})
	if !ret.OK {
		return ret
	}
	// 重算受影响队伍的分数与排名
	if _, _, ret := UpdateTeamRanking(tx, contest, -1, -1); !ret.OK {
		return ret
	}
	return model.SuccessRetVal(map[string]any{
		"cleared_teams":       clearedTeams,
		"cleared_submissions": clearedSubmissions,
	})
}
