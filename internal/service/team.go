package service

import (
	"CBCTF/internal/db"
	"CBCTF/internal/dto"
	"CBCTF/internal/i18n"
	"CBCTF/internal/model"
	"CBCTF/internal/utils"
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
	if int(repo.CountAssociation(team, "Users"))+1 > contest.Size {
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
	contestFlags, ret := db.InitContestFlagRepo(tx).GetTeamSolvedContestFlags(team.ID)
	if !ret.OK {
		return 0, ret
	}
	totalScore := 0.0
	submissionRepo := db.InitSubmissionRepo(tx)
	for _, contestFlag := range contestFlags {
		_, score, ret := CalcContestFlagState(tx, contestFlag)
		if !ret.OK {
			continue
		}
		var rate float64
		if blood {
			bloodTeam, _ := submissionRepo.GetBloodTeamID(contestFlag.ID)
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
		}
		totalScore += score + contestFlag.Score*rate
	}
	totalScore = math.Trunc(totalScore*100) / 100
	return totalScore, model.SuccessRetVal()
}
