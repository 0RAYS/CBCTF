package service

import (
	"CBCTF/internal/db"
	"CBCTF/internal/dto"
	"CBCTF/internal/i18n"
	"CBCTF/internal/model"
	"CBCTF/internal/utils"
	"CBCTF/internal/view"
	"math"
	"sync"
	"time"

	"gorm.io/gorm"
)

func UpdateTeam(tx *gorm.DB, team model.Team, form dto.UpdateTeamForm) model.RetVal {
	repo := db.InitTeamRepo(tx)
	if form.CaptainID != nil && *form.CaptainID != team.CaptainID {
		if !repo.IsInTeam(team.ID, *form.CaptainID) {
			return model.RetVal{Msg: i18n.Model.Team.NotHasMember}
		}
	}
	return repo.Update(team.ID, db.UpdateTeamOptions{
		Description: form.Description,
		Name:        form.Name,
		CaptainID:   form.CaptainID,
	})
}

func BuildTeamView(tx *gorm.DB, team model.Team) view.TeamView {
	count, _ := db.InitTeamRepo(tx).CountUsers(team.ID)
	return view.TeamView{
		Team:      team,
		UserCount: count,
	}
}

func BuildTeamViews(tx *gorm.DB, teams []model.Team) []view.TeamView {
	views := make([]view.TeamView, 0, len(teams))
	if len(teams) == 0 {
		return views
	}

	teamIDs := make([]uint, 0, len(teams))
	for _, team := range teams {
		teamIDs = append(teamIDs, team.ID)
	}
	userCountMap, _ := db.InitTeamRepo(tx).CountUsersMap(teamIDs...)
	for _, team := range teams {
		views = append(views, view.TeamView{
			Team:      team,
			UserCount: userCountMap[team.ID],
		})
	}
	return views
}

func AdminUpdateTeam(tx *gorm.DB, team model.Team, form dto.AdminUpdateTeamForm) model.RetVal {
	repo := db.InitTeamRepo(tx)
	if form.CaptainID != nil && *form.CaptainID != team.CaptainID {
		if !repo.IsInTeam(team.ID, *form.CaptainID) {
			return model.RetVal{Msg: i18n.Model.Team.NotHasMember}
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
	var (
		repo      = db.InitTeamRepo(tx)
		team, ret = repo.GetByName(contest.ID, form.Name)
	)
	if !ret.OK {
		return model.Team{}, ret
	}
	if team.Banned {
		return model.Team{}, model.RetVal{Msg: i18n.Model.Team.Banned}
	}
	if form.Captcha != team.Captcha {
		return model.Team{}, model.RetVal{Msg: i18n.Model.Team.CaptchaWrong}
	}
	mu, _ := JoinTeamMutex.LoadOrStore(team.ID, &sync.Mutex{})
	mu.(*sync.Mutex).Lock()
	defer mu.(*sync.Mutex).Unlock()
	memberCount, ret := repo.CountUsers(team.ID)
	if !ret.OK {
		return model.Team{}, ret
	}
	if int(memberCount)+1 > contest.Size {
		return model.Team{}, model.RetVal{Msg: i18n.Model.Team.Full}
	}
	if repo.IsInContest(contest.ID, user.ID) {
		return model.Team{}, model.RetVal{Msg: i18n.Model.Contest.DuplicateMember}
	}
	if ret = db.AppendUserToTeam(tx, user, team); !ret.OK {
		return model.Team{}, ret
	}
	// 关联 User Contest Many2Many
	if ret = db.AppendUserToContest(tx, user, contest); !ret.OK {
		return model.Team{}, ret
	}
	return team, model.SuccessRetVal()
}

func CreateTeam(tx *gorm.DB, contest model.Contest, user model.User, form dto.CreateTeamForm) (model.Team, model.RetVal) {
	if contest.Captcha != "" && form.Captcha != contest.Captcha {
		return model.Team{}, model.RetVal{Msg: i18n.Model.Contest.CaptchaWrong}
	}
	repo := db.InitTeamRepo(tx)
	if repo.IsInContest(contest.ID, user.ID) {
		return model.Team{}, model.RetVal{Msg: i18n.Model.Contest.DuplicateMember}
	}
	team, ret := repo.Create(db.CreateTeamOptions{
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
		return model.Team{}, ret
	}
	if ret = db.AppendUserToTeam(tx, user, team); !ret.OK {
		return model.Team{}, ret
	}
	if ret = db.AppendUserToContest(tx, user, contest); !ret.OK {
		return model.Team{}, ret
	}
	return team, model.SuccessRetVal()
}

func LeaveTeam(tx *gorm.DB, contest model.Contest, team model.Team, userID uint) model.RetVal {
	repo := db.InitTeamRepo(tx)
	if !repo.IsInTeam(team.ID, userID) {
		return model.RetVal{Msg: i18n.Model.Team.NotHasMember}
	}
	if team.CaptainID == userID {
		return model.RetVal{Msg: i18n.Model.Team.CaptainCannotLeave}
	}
	if ret := db.DeleteUserFromTeam(tx, model.User{BaseModel: model.BaseModel{ID: userID}}, team); !ret.OK {
		return ret
	}
	if ret := db.DeleteUserFromContest(tx, model.User{BaseModel: model.BaseModel{ID: userID}}, contest); !ret.OK {
		return ret
	}
	return model.SuccessRetVal()
}

func CalcTeamScore(tx *gorm.DB, team model.Team, blood bool) (float64, model.RetVal) {
	scoreMap, ret := CalcTeamScores(tx, blood, team)
	if !ret.OK {
		return 0, ret
	}
	score := math.Trunc(scoreMap[team.ID]*100) / 100
	return score, model.SuccessRetVal()
}

func CalcTeamScores(tx *gorm.DB, blood bool, teams ...model.Team) (map[uint]float64, model.RetVal) {
	scoreMap := make(map[uint]float64)
	if len(teams) == 0 {
		return scoreMap, model.SuccessRetVal()
	}

	teamIDL := make([]uint, 0, len(teams))
	for _, team := range teams {
		teamIDL = append(teamIDL, team.ID)
		scoreMap[team.ID] = 0
	}

	rows, ret := db.InitContestFlagRepo(tx).GetTeamsSolvedContestFlags(teamIDL...)
	if !ret.OK {
		return nil, ret
	}
	if len(rows) == 0 {
		return scoreMap, model.SuccessRetVal()
	}

	contestFlagIDL := make([]uint, 0, len(rows))
	for _, row := range rows {
		contestFlagIDL = append(contestFlagIDL, row.ID)
	}

	bloodRanks := make(map[uint]map[uint]int)
	if blood {
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
		scoreMap[row.TeamID] += score
	}

	for teamID, score := range scoreMap {
		scoreMap[teamID] = math.Trunc(score*100) / 100
	}
	return scoreMap, model.SuccessRetVal()
}

func GetTeamView(tx *gorm.DB, team model.Team) view.TeamView {
	return BuildTeamView(tx, team)
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
	return BuildTeamViews(tx, teams), count, model.SuccessRetVal()
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
	return db.InitTeamRepo(tx).Delete(team.ID)
}

func DeleteTeamWithTransaction(tx *gorm.DB, team model.Team) model.RetVal {
	return db.WithTransactionDB(tx, func(tx2 *gorm.DB) model.RetVal {
		return DeleteTeam(tx2, team)
	})
}

func KickMember(tx *gorm.DB, contest model.Contest, team model.Team, userID uint) model.RetVal {
	return db.WithTransactionDB(tx, func(tx2 *gorm.DB) model.RetVal {
		return LeaveTeam(tx2, contest, team, userID)
	})
}

func JoinTeamWithTransaction(tx *gorm.DB, contest model.Contest, user model.User, form dto.JoinTeamForm) (model.Team, model.RetVal) {
	var team model.Team
	ret := db.WithTransactionDB(tx, func(tx2 *gorm.DB) model.RetVal {
		var joinRet model.RetVal
		team, joinRet = JoinTeam(tx2, contest, user, form)
		return joinRet
	})
	return team, ret
}

func CreateTeamWithTransaction(tx *gorm.DB, contest model.Contest, user model.User, form dto.CreateTeamForm) (model.Team, model.RetVal) {
	var team model.Team
	ret := db.WithTransactionDB(tx, func(tx2 *gorm.DB) model.RetVal {
		var createRet model.RetVal
		team, createRet = CreateTeam(tx2, contest, user, form)
		if !createRet.OK {
			return createRet
		}
		return CreateTeamFlags(tx2, team, contest)
	})
	return team, ret
}

func LeaveTeamWithTransaction(tx *gorm.DB, contest model.Contest, team model.Team, userID uint) model.RetVal {
	return db.WithTransactionDB(tx, func(tx2 *gorm.DB) model.RetVal {
		return LeaveTeam(tx2, contest, team, userID)
	})
}
