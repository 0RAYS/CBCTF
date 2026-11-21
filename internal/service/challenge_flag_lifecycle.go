package service

import (
	"fmt"

	"gorm.io/gorm"

	"CBCTF/internal/db"
	"CBCTF/internal/dto"
	"CBCTF/internal/model"
)

// UpdateChallengeFlag 修改题目的 flag，默认保留已解出队伍的解题记录。
// 改 flag 模板时代次 +1，未解出队伍的 team_flags 按新模板重新渲染。
func UpdateChallengeFlag(tx *gorm.DB, flag model.ChallengeFlag, form dto.UpdateChallengeFlagForm) model.RetVal {
	currentScore := flag.CurrentScore
	if form.Score != nil && *form.Score < currentScore {
		currentScore = *form.Score
	}
	options := db.UpdateChallengeFlagOptions{
		Score:        form.Score,
		CurrentScore: &currentScore,
		Decay:        form.Decay,
		MinScore:     form.MinScore,
		ScoreType:    form.ScoreType,
	}
	if form.Value == nil || *form.Value == flag.Value {
		return db.InitChallengeFlagRepo(tx).Update(flag.ID, options)
	}
	generation := flag.Generation + 1
	options.Value = form.Value
	options.Generation = &generation
	if ret := db.InitChallengeFlagRepo(tx).Update(flag.ID, options); !ret.OK {
		return ret
	}
	contest, ret := challengeContest(tx, flag.ChallengeID)
	if !ret.OK {
		return ret
	}
	return RerenderTeamFlagValues(tx, contest, flag.ID, *form.Value)
}

// challengeContest 返回题目所属的比赛，题目不归属比赛时返回零值。
func challengeContest(tx *gorm.DB, challengeID uint) (model.Contest, model.RetVal) {
	challenge, ret := db.InitChallengeRepo(tx).GetByID(challengeID)
	if !ret.OK {
		return model.Contest{}, ret
	}
	if !challenge.ContestID.Valid {
		return model.Contest{}, model.SuccessRetVal()
	}
	return db.InitContestRepo(tx).GetByID(challenge.ContestID.V)
}

// RerenderTeamFlagValues 按新的 flag 模板重新渲染未解出队伍手里的答案，已解出的队伍保留原值。
func RerenderTeamFlagValues(tx *gorm.DB, contest model.Contest, challengeFlagID uint, template string) model.RetVal {
	teamFlags, _, ret := db.InitTeamFlagRepo(tx).List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"challenge_flag_id": challengeFlagID},
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

// ResetChallengeSolves 清空某题目的解题记录并重算分数与排名。
func ResetChallengeSolves(tx *gorm.DB, challenge model.Challenge) model.RetVal {
	contest, ret := challengeContest(tx, challenge.ID)
	if !ret.OK {
		return ret
	}
	var clearedTeams, clearedSubmissions int
	ret = db.WithTransactionDB(tx, func(tx2 *gorm.DB) model.RetVal {
		flags, _, ret := db.InitChallengeFlagRepo(tx2).List(-1, -1, db.GetOptions{
			Conditions: map[string]any{"challenge_id": challenge.ID},
		})
		if !ret.OK {
			return ret
		}
		flagIDL := make([]uint, 0, len(flags))
		for _, flag := range flags {
			flagIDL = append(flagIDL, flag.ID)
		}
		if len(flagIDL) > 0 {
			teamFlags, _, ret := db.InitTeamFlagRepo(tx2).List(-1, -1, db.GetOptions{
				Conditions: map[string]any{"challenge_flag_id": flagIDL},
			})
			if !ret.OK {
				return ret
			}
			teams := make(map[uint]struct{}, len(teamFlags))
			for _, teamFlag := range teamFlags {
				teams[teamFlag.TeamID] = struct{}{}
			}
			clearedTeams = len(teams)
			if ret := db.InitTeamFlagRepo(tx2).DeleteByChallengeFlagID(flagIDL...); !ret.OK {
				return ret
			}
		}
		solvedSubmissions, _, ret := db.InitSubmissionRepo(tx2).List(-1, -1, db.GetOptions{
			Conditions: map[string]any{
				"challenge_id": challenge.ID,
				"solved":       true,
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
		flagRepo := db.InitChallengeFlagRepo(tx2)
		for _, flag := range flags {
			score := flag.Score
			if ret := flagRepo.Update(flag.ID, db.UpdateChallengeFlagOptions{
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
	if contest.ID != 0 {
		if _, _, ret := UpdateTeamRanking(tx, contest, -1, -1); !ret.OK {
			return ret
		}
	}
	return model.SuccessRetVal(map[string]any{
		"cleared_teams":       clearedTeams,
		"cleared_submissions": clearedSubmissions,
	})
}
