package service

import (
	"CBCTF/internal/db"
	f "CBCTF/internal/form"
	"CBCTF/internal/i18n"
	"CBCTF/internal/model"
	"CBCTF/internal/prometheus"
	"CBCTF/internal/utils"
	"math"
	"time"

	"gorm.io/gorm"
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

func JoinTeam(tx *gorm.DB, contest model.Contest, user model.User, form f.JoinTeamForm) (model.Team, bool, string) {
	var (
		repo          = db.InitTeamRepo(tx)
		team, ok, msg = repo.GetByName(contest.ID, form.Name, db.GetOptions{
			Preloads: map[string]db.GetOptions{"Users": {}},
		})
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
	prometheus.AddContestActiveUsersMetrics(contest, 1)
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
	if repo.IsInContest(contest.ID, user.ID) {
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
	prometheus.AddContestActiveTeamsMetrics(contest, 1)
	prometheus.AddContestActiveUsersMetrics(contest, 1)
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
	prometheus.SubContestActiveUsersMetrics(contest, 1)
	return true, i18n.Success
}

func GetTeamSolvedFlags(tx *gorm.DB, team model.Team) ([]model.ContestFlag, bool, string) {
	solvedContestFlags := make([]model.ContestFlag, 0)
	solvedSubmissions, _, ok, msg := db.InitSubmissionRepo(tx).List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"team_id": team.ID, "solved": true},
		Preloads:   map[string]db.GetOptions{"ContestFlag": {}},
	})
	if !ok {
		return nil, false, msg
	}
	for _, submission := range solvedSubmissions {
		solvedContestFlags = append(solvedContestFlags, submission.ContestFlag)
	}
	return solvedContestFlags, true, i18n.Success
}

func CalcTeamScore(tx *gorm.DB, team model.Team) (float64, bool, string) {
	submissionRepo := db.InitSubmissionRepo(tx)
	submissions, _, ok, msg := submissionRepo.List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"team_id": team.ID, "solved": true},
		Preloads:   map[string]db.GetOptions{"ContestFlag": {}},
	})
	if !ok {
		return 0, false, msg
	}
	totalScore := 0.0
	for _, submission := range submissions {
		_, score, ok, _ := CalcContestFlagState(tx, submission.ContestFlag)
		if !ok {
			continue
		}
		var rate float64
		bloodTeam, _, _ := submissionRepo.GetBloodTeam(submission.ContestFlagID)
		for i, teamID := range bloodTeam {
			if teamID == team.ID {
				switch i {
				case 0:
					rate = model.FirstBloodRate
				case 1:
					rate = model.SecondBloodRate
				case 2:
					rate = model.ThirdBloodRate
				}
			}
			if rate > 0 {
				break
			}
		}
		totalScore += score + submission.ContestFlag.Score*rate
	}
	totalScore = math.Trunc(totalScore*100) / 100
	return totalScore, true, i18n.Success
}
