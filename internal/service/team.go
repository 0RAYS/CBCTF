package service

import (
	"fmt"
	"math"
	"sync"
	"time"

	"gorm.io/gorm"

	"CBCTF/internal/db"
	"CBCTF/internal/dto"
	"CBCTF/internal/i18n"
	"CBCTF/internal/model"
	"CBCTF/internal/utils"
	"CBCTF/internal/view"
)

func UpdateTeam(tx *gorm.DB, team model.Team, form dto.UpdateTeamForm) model.RetVal {
	repo := db.InitTeamRepo(tx)
	if form.CaptainID != nil && *form.CaptainID != team.CaptainID {
		if !repo.IsInTeam(team.ID, *form.CaptainID) {
			return model.FailRetVal(i18n.Model.Team.NotHasMember, "Team", fmt.Sprintf("team_id=%d,user_id=%d", team.ID, *form.CaptainID), nil)
		}
	}
	return repo.Update(team.ID, db.UpdateTeamOptions{
		Description: form.Description,
		Name:        form.Name,
		CaptainID:   form.CaptainID,
	})
}

func BuildTeamView(tx *gorm.DB, team model.Team) view.TeamView {
	count, ret := db.InitTeamRepo(tx).CountUsers(team.ID)
	result := view.TeamView{
		Team:      team,
		UserCount: count,
	}
	if !ret.OK {
		result.Unavailable = []string{"users"}
	}
	adjustments, _ := db.InitScoreAdjustmentRepo(tx).SumByTeamID(team.ContestID, team.ID)
	result.Adjustment = adjustments[team.ID]
	return result
}

func AdminUpdateTeam(tx *gorm.DB, team model.Team, form dto.AdminUpdateTeamForm) model.RetVal {
	repo := db.InitTeamRepo(tx)
	if form.CaptainID != nil && *form.CaptainID != team.CaptainID {
		if !repo.IsInTeam(team.ID, *form.CaptainID) {
			return model.FailRetVal(i18n.Model.Team.NotHasMember, "Team", fmt.Sprintf("team_id=%d,user_id=%d", team.ID, *form.CaptainID), nil)
		}
	}
	return repo.Update(team.ID, db.UpdateTeamOptions{
		Name:        form.Name,
		Description: form.Description,
		Hidden:      form.Hidden,
		Banned:      form.Banned,
		Captcha:     form.Captcha,
		CaptainID:   form.CaptainID,
	})
}

var JoinTeamMutex sync.Map

func JoinTeam(tx *gorm.DB, contest model.Contest, user model.User, form dto.JoinTeamForm) (model.Team, model.RetVal) {
	var result model.Team
	ret := db.WithTransactionDB(tx, func(tx *gorm.DB) model.RetVal {
		var (
			repo      = db.InitTeamRepo(tx)
			team, ret = repo.GetByName(contest.ID, form.Name)
		)
		if !ret.OK {
			return ret
		}
		if team.Banned {
			return model.FailRetVal(i18n.Model.Team.Banned, "Team", fmt.Sprintf("team_id=%d", team.ID), nil)
		}
		if form.Captcha != team.Captcha {
			return model.FailRetVal(i18n.Model.Team.CaptchaWrong, "Team", fmt.Sprintf("team_id=%d", team.ID), nil)
		}
		mu, _ := JoinTeamMutex.LoadOrStore(team.ID, &sync.Mutex{})
		mu.(*sync.Mutex).Lock()
		defer mu.(*sync.Mutex).Unlock()
		memberCount, ret := repo.CountUsers(team.ID)
		if !ret.OK {
			return ret
		}
		if int(memberCount)+1 > contest.Size {
			return model.FailRetVal(i18n.Model.Team.Full, "Team", fmt.Sprintf("team_id=%d,contest_id=%d", team.ID, contest.ID), nil)
		}
		if repo.IsInContest(contest.ID, user.ID) {
			return model.FailRetVal(i18n.Model.Contest.DuplicateMember, "Contest", fmt.Sprintf("contest_id=%d,user_id=%d", contest.ID, user.ID), nil)
		}
		if ret = db.AppendUserToTeam(tx, user, team); !ret.OK {
			return ret
		}
		if ret = db.AppendUserToContest(tx, user, contest); !ret.OK {
			return ret
		}
		result = team
		return model.SuccessRetVal()
	})
	return result, ret
}

func CreateTeam(tx *gorm.DB, contest model.Contest, user model.User, form dto.CreateTeamForm) (model.Team, model.RetVal) {
	var result model.Team
	ret := db.WithTransactionDB(tx, func(tx *gorm.DB) model.RetVal {
		if contest.Captcha != "" && form.Captcha != contest.Captcha {
			return model.FailRetVal(i18n.Model.Contest.CaptchaWrong, "Contest", fmt.Sprintf("contest_id=%d", contest.ID), nil)
		}
		repo := db.InitTeamRepo(tx)
		if repo.IsInContest(contest.ID, user.ID) {
			return model.FailRetVal(i18n.Model.Contest.DuplicateMember, "Contest", fmt.Sprintf("contest_id=%d,user_id=%d", contest.ID, user.ID), nil)
		}
		team, ret := repo.Create(model.Team{
			Name:        form.Name,
			ContestID:   contest.ID,
			Description: form.Description,
			Captcha:     utils.UUID(),
			Picture:     "",
			Banned:      false,
			Hidden:      false,
			CaptainID:   user.ID,
			Last:        time.Now(),
		})
		if !ret.OK {
			return ret
		}
		if ret = db.AppendUserToTeam(tx, user, team); !ret.OK {
			return ret
		}
		if ret = db.AppendUserToContest(tx, user, contest); !ret.OK {
			return ret
		}
		result = team
		return CreateTeamFlags(tx, team, contest)
	})
	return result, ret
}

func LeaveTeam(tx *gorm.DB, contest model.Contest, team model.Team, userID uint) model.RetVal {
	return db.WithTransactionDB(tx, func(tx *gorm.DB) model.RetVal {
		repo := db.InitTeamRepo(tx)
		if !repo.IsInTeam(team.ID, userID) {
			return model.FailRetVal(i18n.Model.Team.NotHasMember, "Team", fmt.Sprintf("team_id=%d,user_id=%d", team.ID, userID), nil)
		}
		if team.CaptainID == userID {
			return model.FailRetVal(i18n.Model.Team.CaptainCannotLeave, "Team", fmt.Sprintf("team_id=%d,user_id=%d", team.ID, userID), nil)
		}
		if ret := db.DeleteUserFromTeam(tx, model.User{ID: userID}, team); !ret.OK {
			return ret
		}
		if ret := db.DeleteUserFromContest(tx, model.User{ID: userID}, contest); !ret.OK {
			return ret
		}
		return model.SuccessRetVal()
	})
}

func CalcTeamScore(tx *gorm.DB, contest model.Contest, team model.Team) (float64, model.RetVal) {
	scoreMap, ret := CalcTeamScores(tx, contest.ID, contest.Blood, team)
	if !ret.OK {
		return 0, ret
	}
	return scoreMap[team.ID].Total, model.SuccessRetVal()
}

// TeamScore 队伍总分及其分题得分明细。
// ByChallenge 的键即该队伍解出的比赛题目，未解出的题目不会出现在其中。
type TeamScore struct {
	Total       float64
	ByChallenge map[uint]float64
}

// CalcTeamScores 计算队伍总分。总分等于解题分数加上管理员对队伍的分数变动总和，
// 变动分数只计入总分，不计入分题得分。
func CalcTeamScores(tx *gorm.DB, contestID uint, blood bool, teams ...model.Team) (map[uint]TeamScore, model.RetVal) {
	scoreMap := make(map[uint]TeamScore)
	if len(teams) == 0 {
		return scoreMap, model.SuccessRetVal()
	}

	teamIDL := make([]uint, 0, len(teams))
	for _, team := range teams {
		teamIDL = append(teamIDL, team.ID)
		scoreMap[team.ID] = TeamScore{ByChallenge: make(map[uint]float64)}
	}

	adjustments, ret := db.InitScoreAdjustmentRepo(tx).SumByTeamID(contestID, teamIDL...)
	if !ret.OK {
		return nil, ret
	}

	rows, ret := db.InitContestFlagRepo(tx).GetTeamsSolvedContestFlags(teamIDL...)
	if !ret.OK {
		return nil, ret
	}

	contestFlagIDL := make([]uint, 0, len(rows))
	for _, row := range rows {
		contestFlagIDL = append(contestFlagIDL, row.ID)
	}

	bloodRanks := make(map[uint]map[uint]int)
	if blood && len(contestFlagIDL) > 0 {
		bloodRanks, ret = db.InitSubmissionRepo(tx).GetBloodRankMap(contestFlagIDL...)
		if !ret.OK {
			return nil, ret
		}
	}

	for _, row := range rows {
		score := row.CurrentScore
		if blood {
			switch bloodRanks[row.ID][row.TeamID] {
			case 1:
				score += row.Score * model.FirstBloodRate
			case 2:
				score += row.Score * model.SecondBloodRate
			case 3:
				score += row.Score * model.ThirdBloodRate
			}
		}
		teamScore := scoreMap[row.TeamID]
		teamScore.Total += score
		teamScore.ByChallenge[row.ContestChallengeID] += score
		scoreMap[row.TeamID] = teamScore
	}

	finalizeTeamScores(scoreMap, adjustments)
	return scoreMap, model.SuccessRetVal()
}

// finalizeTeamScores 规整队伍分数，并把管理员对队伍的分数变动计入总分。分题得分不受变动影响。
func finalizeTeamScores(scoreMap map[uint]TeamScore, adjustments map[uint]float64) {
	for teamID, teamScore := range scoreMap {
		teamScore.Total = math.Trunc((teamScore.Total+adjustments[teamID])*100) / 100
		for challengeID, score := range teamScore.ByChallenge {
			teamScore.ByChallenge[challengeID] = math.Trunc(score*100) / 100
		}
		scoreMap[teamID] = teamScore
	}
}

func ListTeams(tx *gorm.DB, contest model.Contest, form dto.ListTeamForm) ([]view.TeamView, int64, model.RetVal) {
	options := db.GetOptions{
		Conditions: map[string]any{"contest_id": contest.ID},
		Search:     make(map[string]string),
	}
	if form.Name != "" {
		options.Search["name"] = form.Name
	}
	if form.Description != "" {
		options.Search["description"] = form.Description
	}
	teams, count, ret := db.InitTeamRepo(tx).List(form.Limit, form.Offset, options)
	if !ret.OK {
		return nil, 0, ret
	}
	views := make([]view.TeamView, 0, len(teams))
	if len(teams) == 0 {
		return views, count, model.SuccessRetVal()
	}
	teamIDs := make([]uint, 0, len(teams))
	for _, team := range teams {
		teamIDs = append(teamIDs, team.ID)
	}
	userCountMap, _ := db.InitTeamRepo(tx).CountUsersMap(teamIDs...)
	adjustments, _ := db.InitScoreAdjustmentRepo(tx).SumByTeamID(contest.ID, teamIDs...)
	for _, team := range teams {
		views = append(views, view.TeamView{
			Team:       team,
			UserCount:  userCountMap[team.ID],
			Adjustment: adjustments[team.ID],
		})
	}
	return views, count, model.SuccessRetVal()
}

func GetTeammates(tx *gorm.DB, team model.Team, includeCounts bool) ([]view.UserView, model.RetVal) {
	users, ret := db.InitUserRepo(tx).GetByTeamID(team.ID, -1, -1)
	if !ret.OK {
		return nil, ret
	}
	return BuildUserViews(tx, users, includeCounts), model.SuccessRetVal()
}

func GetTeamSolvedFlags(tx *gorm.DB, contest model.Contest, team model.Team) ([]model.ContestFlag, []model.ContestFlag, model.RetVal) {
	contestFlagRepo := db.InitContestFlagRepo(tx)
	contestFlagL, _, ret := contestFlagRepo.List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"contest_id": contest.ID},
		Preloads:   map[string]db.GetOptions{"ContestChallenge": {}},
	})
	if !ret.OK {
		return nil, nil, ret
	}
	solvedFlagL, _ := contestFlagRepo.GetTeamSolvedContestFlags(team.ID)
	return solvedFlagL, contestFlagL, model.SuccessRetVal()
}

func UpdateTeamCaptcha(tx *gorm.DB, team model.Team, captcha string) model.RetVal {
	return db.InitTeamRepo(tx).Update(team.ID, db.UpdateTeamOptions{Captcha: &captcha})
}

func DeleteTeam(tx *gorm.DB, team model.Team) model.RetVal {
	return db.WithTransactionDB(tx, func(tx *gorm.DB) model.RetVal {
		return db.InitTeamRepo(tx).Delete(team.ID)
	})
}
