package service

import (
	f "CBCTF/internel/form"
	"CBCTF/internel/model"
	db "CBCTF/internel/repo"
	"CBCTF/internel/utils"
	"gorm.io/gorm"
	"math"
)

func UpdateTeam(tx *gorm.DB, team model.Team, form f.UpdateTeamForm) (bool, string) {
	repo := db.InitTeamRepo(tx)
	if form.Name != nil && *form.Name != team.Name {
		if !repo.IsUniqueName(team.ContestID, *form.Name) {
			return false, "DuplicateTeamName"
		}
	}
	if form.CaptainID != nil && *form.CaptainID != team.CaptainID {
		if !repo.IsTeamMember(team.ID, *form.CaptainID) {
			return false, "UserNotInTeam"
		}
	}
	return repo.Update(team.ID, db.UpdateTeamOptions{
		Desc:      form.Desc,
		Name:      form.Name,
		CaptainID: form.CaptainID,
	})
}

func AdminUpdateTeam(tx *gorm.DB, team model.Team, form f.AdminUpdateTeamForm) (bool, string) {
	repo := db.InitTeamRepo(tx)
	if form.Name != nil && *form.Name != team.Name {
		if !repo.IsUniqueName(team.ContestID, *form.Name) {
			return false, "DuplicateTeamName"
		}
	}
	if form.CaptainID != nil && *form.CaptainID != team.CaptainID {
		if !repo.IsTeamMember(team.ID, *form.CaptainID) {
			return false, "UserNotInTeam"
		}
	}
	return repo.Update(team.ID, db.UpdateTeamOptions{
		Name:      form.Name,
		Desc:      form.Desc,
		Hidden:    form.Hidden,
		Banned:    form.Banned,
		Captcha:   form.Captcha,
		CaptainID: form.CaptainID,
	})
}

func UpdateTeamCaptcha(tx *gorm.DB, team model.Team, captcha string) (bool, string) {
	repo := db.InitTeamRepo(tx)
	return repo.Update(team.ID, db.UpdateTeamOptions{Captcha: &captcha})
}

func DeleteTeam(tx *gorm.DB, team model.Team) (bool, string) {
	repo := db.InitTeamRepo(tx)
	return repo.Delete(team.ID)
}

func JoinTeam(tx *gorm.DB, contest model.Contest, user model.User, form f.JoinTeamForm) (bool, string) {
	var (
		repo          = db.InitTeamRepo(tx)
		team, ok, msg = repo.GetByName(contest.ID, form.Name, true)
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
	if err = db.AppendUserToTeam(tx, user.ID, team.ID); err != nil {
		return false, "AppendUserToTeamError"
	}
	// 关联 User Contest Many2Many
	if err = db.AppendUserToContest(tx, user.ID, contest.ID); err != nil {
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
	if err := db.AppendUserToTeam(tx, user.ID, team.ID); err != nil {
		return false, "AppendUserToTeamError"
	}
	if err := db.AppendUserToContest(tx, user.ID, contest.ID); err != nil {
		return false, "AppendUserToContestError"
	}
	return true, "Success"
}

func LeaveTeam(tx *gorm.DB, contest model.Contest, team model.Team, userID uint) (bool, string) {
	repo := db.InitTeamRepo(tx)
	if !repo.IsTeamMember(team.ID, userID) {
		return false, "UserNotInTeam"
	}
	if team.CaptainID == userID {
		return false, "CaptainCannotLeave"
	}
	if err := db.DeleteUserFromTeam(tx, userID, team.ID); err != nil {
		return false, "DeleteUserFromTeamError"
	}
	if err := db.DeleteUserFromContest(tx, userID, contest.ID); err != nil {
		return false, "DeleteUserFromContestError"
	}
	return true, "Success"
}

func CalcTeamScore(tx *gorm.DB, teamID uint) (float64, bool, string) {
	var (
		teamRepo      = db.InitTeamRepo(tx)
		usageRepo     = db.InitUsageRepo(tx)
		team, ok, msg = teamRepo.GetByID(teamID, true)
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
		usage, ok, msg = usageRepo.GetBy2ID(submission.ContestID, submission.ChallengeID, false, true)
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
		submissions, _, ok, msg = repo.GetAllByKeyID("team_id", teamID, -1, -1, true, true)
	)
	if !ok {
		return flags, false, msg
	}
	for _, submission := range submissions {
		flags = append(flags, submission.Flag)
	}
	return flags, true, "Success"
}
