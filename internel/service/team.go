package service

import (
	f "CBCTF/internel/form"
	"CBCTF/internel/i18n"
	"CBCTF/internel/model"
	db "CBCTF/internel/repo"
	"CBCTF/internel/utils"
	"gorm.io/gorm"
	"time"
)

func UpdateTeam(tx *gorm.DB, team model.Team, form f.UpdateTeamForm) (bool, string) {
	repo := db.InitTeamRepo(tx)
	if form.Name != nil && *form.Name != team.Name {
		if !repo.IsUniqueName(team.ContestID, *form.Name) {
			return false, i18n.DuplicateTeamName
		}
	}
	if form.CaptainID != nil && *form.CaptainID != team.CaptainID {
		if !repo.IsInTeam(team.ID, *form.CaptainID) {
			return false, i18n.UserNotInTeam
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
			return false, i18n.DuplicateTeamName
		}
	}
	if form.CaptainID != nil && *form.CaptainID != team.CaptainID {
		if !repo.IsInTeam(team.ID, *form.CaptainID) {
			return false, i18n.UserNotInTeam
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

func JoinTeam(tx *gorm.DB, contest model.Contest, user model.User, form f.JoinTeamForm) (model.Team, bool, string) {
	var (
		repo          = db.InitTeamRepo(tx)
		team, ok, msg = repo.GetByName(contest.ID, form.Name, "Users")
	)
	if !ok {
		return model.Team{}, false, msg
	}
	if team.Banned {
		return model.Team{}, false, i18n.TeamIsBanned
	}
	if form.Captcha != team.Captcha {
		return model.Team{}, false, i18n.TeamCaptchaError
	}
	if len(team.Users)+1 > contest.Size {
		return model.Team{}, false, i18n.TeamIsFull
	}
	if repo.IsInContest(contest.ID, user.ID) {
		return model.Team{}, false, i18n.DuplicateMember
	}
	if ok, msg = db.AppendUserToTeam(tx, user.ID, team.ID); !ok {
		return model.Team{}, false, msg
	}
	// 关联 User Contest Many2Many
	if ok, msg = db.AppendUserToContest(tx, user.ID, contest.ID); !ok {
		return model.Team{}, false, msg
	}
	team.Users = append(team.Users, &user)
	return team, true, i18n.Success
}

func CreateTeam(tx *gorm.DB, contest model.Contest, user model.User, form f.CreateTeamForm) (model.Team, bool, string) {
	if contest.Captcha != "" && form.Captcha != contest.Captcha {
		return model.Team{}, false, i18n.ContestCaptchaError
	}
	repo := db.InitTeamRepo(tx)
	if !repo.IsUniqueName(contest.ID, form.Name) {
		return model.Team{}, false, i18n.DuplicateTeamName
	}
	if repo.IsInTeam(contest.ID, user.ID) {
		return model.Team{}, false, i18n.DuplicateMember
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
		Last:      time.Now(),
	})
	if !ok {
		return model.Team{}, false, msg
	}
	if ok, msg = db.AppendUserToTeam(tx, user.ID, team.ID); !ok {
		return model.Team{}, false, msg
	}
	if ok, msg = db.AppendUserToContest(tx, user.ID, contest.ID); !ok {
		return model.Team{}, false, msg
	}
	team.Users = append(team.Users, &user)
	return team, true, i18n.Success
}

func LeaveTeam(tx *gorm.DB, contest model.Contest, team model.Team, userID uint) (bool, string) {
	repo := db.InitTeamRepo(tx)
	if !repo.IsInTeam(team.ID, userID) {
		return false, i18n.UserNotInTeam
	}
	if team.CaptainID == userID {
		return false, i18n.CaptainCannotLeave
	}
	if ok, msg := db.DeleteUserFromTeam(tx, userID, team.ID); !ok {
		return false, msg
	}
	if ok, msg := db.DeleteUserFromContest(tx, userID, contest.ID); !ok {
		return false, msg
	}
	return true, i18n.Success
}

// GetTeamSolved TODO
//func GetTeamSolved(tx *gorm.DB, teamID uint) ([]model.Flag, bool, string) {
//var (
//	flags                   = make([]model.Flag, 0)
//	repo                    = db.InitSubmissionRepo(tx)
//	submissions, _, ok, msg = repo.GetByKeyID("team_id", teamID, -1, -1, true, "Flag")
//)
//if !ok {
//	return flags, false, msg
//}
//for _, submission := range submissions {
//	flags = append(flags, submission.Flag)
//}
//return flags, true, i18n.Success
//}
