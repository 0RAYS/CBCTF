package service

import (
	f "CBCTF/internel/form"
	"CBCTF/internel/model"
	"CBCTF/internel/redis"
	db "CBCTF/internel/repo"
	"CBCTF/internel/utils"
	"gorm.io/gorm"
	"math"
)

func UpdateTeamCaptcha(tx *gorm.DB, team model.Team, captcha string) (bool, string) {
	repo := db.InitTeamRepo(tx)
	return repo.Update(team.ID, db.UpdateTeamOptions{Captcha: &captcha})
}

func JoinTeam(tx *gorm.DB, contest model.Contest, user model.User, form f.JoinTeamForm) (bool, string) {
	var (
		repo          = db.InitTeamRepo(tx)
		team, ok, msg = repo.GetByName(contest.ID, form.Name, true, 0)
		err           error
	)
	if !ok {
		return false, msg
	}
	if team.Banned {
		return false, "TeamBanned"
	}
	if form.Captcha != team.Captcha {
		return false, "CaptchaError"
	}
	if len(team.Users)+1 > contest.Size {
		return false, "TeamIsFull"
	}
	if !repo.IsUniqueMember(contest.ID, user.ID) {
		return false, "DuplicateMember"
	}
	if err = db.AppendUserToTeam(tx, user, team); err != nil {
		return false, "AppendUserToTeamError"
	}
	// 关联 User Contest Many2Many
	if err = db.AppendUserToContest(tx, user, contest); err != nil {
		return false, "AppendContestToUserError"
	}
	return true, "Success"
}

func CreateTeam(tx *gorm.DB, contest model.Contest, user model.User, form f.CreateTeamForm) (bool, string) {
	if contest.Captcha != "" && form.Captcha != contest.Captcha {
		return false, "CaptchaError"
	}
	repo := db.InitTeamRepo(tx)
	if !repo.IsUniqueName(contest.ID, form.Name) {
		return false, "DuplicateTeamName"
	}
	if !repo.IsUniqueMember(contest.ID, user.ID) {
		return false, "DuplicateMember"
	}
	team, ok, msg := repo.Create(db.CreateTeamOptions{
		Name:      form.Name,
		ContestID: contest.ID,
		Desc:      form.Desc,
		Captcha:   utils.UUID(),
		Avatar:    "",
		Banned:    false,
		Hidden:    false,
		CaptainID: user.ID,
	})
	if !ok {
		return false, msg
	}
	if err := db.AppendUserToTeam(tx, user, team); err != nil {
		return false, "AppendUserToTeamError"
	}
	if err := db.AppendUserToContest(tx, user, contest); err != nil {
		return false, "AppendUserToContestError"
	}
	return true, "Success"
}

func UpdateTeamRanking(tx *gorm.DB, contestID uint) (bool, string) {
	var (
		repo              = db.InitTeamRepo(tx)
		teams, _, ok, msg = repo.GetAll(contestID, -1, -1, true, 0, true, false)
		score             float64
		err               error
	)
	if !ok {
		return false, msg
	}
	for _, team := range teams {
		score, ok, msg = CalcTeamScore(tx, team.ID)
		if !ok {
			continue
		}
		// 不考虑更新失败的情况, 不回滚
		repo.Update(team.ID, db.UpdateTeamOptions{Score: &score})
	}
	teams, _, ok, msg = repo.GetAll(contestID, -1, -1, true, 0, true, false)
	if !ok {
		return false, msg
	}
	if err = redis.UpdateTeamRanking(contestID, teams); err != nil {
		return false, "UpdateRankingError"
	}
	return true, "Success"
}

func GetTeamRanking(tx *gorm.DB, contestID uint, limit, offset int) ([]model.Team, int64, bool, string) {
	var (
		teams          = make([]model.Team, 0)
		repo           = db.InitTeamRepo(tx)
		count, ok, msg = repo.Count(contestID, true, false)
		err            error
	)
	if !ok {
		return teams, count, false, msg
	}
	start, end := utils.TidyPaginate(int(count), limit, offset)
	if teams, err = redis.GetTeamRanking(contestID, int64(start), int64(end-1)); err == nil && teams != nil {
		return teams, count, true, "Success"
	}
	if ok, msg = UpdateTeamRanking(tx, contestID); err != nil {
		return teams, count, false, msg
	}
	return GetTeamRanking(tx, contestID, limit, offset)
}

func CalcTeamScore(tx *gorm.DB, teamID uint) (float64, bool, string) {
	var (
		teamRepo      = db.InitTeamRepo(tx)
		usageRepo     = db.InitUsageRepo(tx)
		team, ok, msg = teamRepo.GetByID(teamID, true, 3)
		usage         model.Usage
		total         float64
		score         float64
	)
	if !ok {
		return team.Score, false, msg
	}
	for _, submission := range team.Submissions {
		if !submission.Solved {
			continue
		}
		usage, ok, msg = usageRepo.GetBy2ID(submission.ContestID, submission.ChallengeID, true, 3, false)
		if !ok {
			continue
		}
		for _, flag := range usage.Flags {
			_, score, ok, msg = CalcSolversAndScore(tx, flag)
			if !ok {
				continue
			}
			rate, _ := flag.CalcBlood(team.ID)
			total += score + flag.Score*rate
		}
	}
	score = math.Trunc(score*100) / 100
	return total, true, "Success"
}

func GetTeamSolved(tx *gorm.DB, teamID uint) ([]model.Flag, bool, string) {
	var (
		flags                   = make([]model.Flag, 0)
		repo                    = db.InitSubmissionRepo(tx)
		submissions, _, ok, msg = repo.GetAllByKeyID("team_id", teamID, -1, -1, true, 0, true)
	)
	if !ok {
		return flags, false, msg
	}
	for _, submission := range submissions {
		flags = append(flags, submission.Flag)
	}
	return flags, true, "Success"
}
