package repo

import (
	"CBCTF/internel/i18n"
	"CBCTF/internel/log"
	"CBCTF/internel/model"
	"CBCTF/internel/utils"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type TeamRepo struct {
	Repo[model.Team]
}

type CreateTeamOptions struct {
	Name      string
	ContestID uint
	Desc      string
	Captcha   string
	Avatar    string
	Banned    bool
	Hidden    bool
	CaptainID uint
	Last      time.Time
}

type UpdateTeamOptions struct {
	Name      *string    `json:"name"`
	Desc      *string    `json:"desc"`
	Captcha   *string    `json:"captcha"`
	Avatar    *string    `json:"avatar"`
	Banned    *bool      `json:"banned"`
	Hidden    *bool      `json:"hidden"`
	CaptainID *uint      `json:"captain_id"`
	Score     *float64   `json:"score"`
	Rank      *int       `json:"rank"`
	Last      *time.Time `json:"last"`
}

func InitTeamRepo(tx *gorm.DB) *TeamRepo {
	return &TeamRepo{Repo: Repo[model.Team]{DB: tx, Model: "Team"}}
}

func (t *TeamRepo) IsUniqueName(contestID uint, name string) bool {
	res := t.DB.Model(&model.Team{}).Where("contest_id = ? AND name = ?", contestID, name).Limit(1).
		Find(&model.Team{})
	return res.RowsAffected == 0
}

func (t *TeamRepo) IsTeamMember(teamID uint, userID uint) bool {
	res := t.DB.Model(&model.UserTeam{}).
		Where("team_id = ? AND user_id = ?", teamID, userID).Limit(1).Find(&model.UserTeam{})
	return res.RowsAffected == 1
}

func (t *TeamRepo) IsUniqueMember(contestID uint, userID uint) bool {
	res := t.DB.Model(&model.UserContest{}).
		Where("contest_id = ? AND user_id = ?", contestID, userID).Limit(1).Find(&model.UserContest{})
	return res.RowsAffected == 0
}

func (t *TeamRepo) GetByName(contestID uint, name string, preloadL ...string) (model.Team, bool, string) {
	var team model.Team
	res := t.DB.Model(&model.Team{}).Where("contest_id = ? AND name = ?", contestID, name)
	res = preload(res, preloadL...).Limit(1).Find(&team)
	if res.Error != nil {
		log.Logger.Warningf("Failed to get Team")
		return model.Team{}, false, i18n.GetTeamError
	}
	if res.RowsAffected == 0 {
		return model.Team{}, false, i18n.TeamNotFound
	}
	return team, true, i18n.Success
}

// GetBy2ID 根据用户 ID 和比赛 ID 获取 model.Team, 等同于 GetByID(teamID, "all")
func (t *TeamRepo) GetBy2ID(userID uint, contestID uint) (model.Team, bool, string) {
	user, ok, msg := InitUserRepo(t.DB).
		GetByID(
			userID,
			"all",
			"Teams.Contest", "Teams.Users", "Teams.Answers", "Teams.Submissions",
			"Teams.Victims", "Teams.Cheats",
		)
	if !ok {
		return model.Team{}, false, msg
	}
	for _, team := range user.Teams {
		if team.ContestID == contestID {
			return *team, true, i18n.Success
		}
	}
	return model.Team{}, false, i18n.UserNotInTeam
}

func (t *TeamRepo) Count(contestID uint, hidden, banned bool) (int64, bool, string) {
	var count int64
	res := t.DB.Model(&model.Team{}).Where("contest_id = ?", contestID)
	if !hidden {
		res = res.Where("hidden = ?", false)
	}
	if !banned {
		res = res.Where("banned = ?", false)
	}
	res = res.Count(&count)
	if res.Error != nil {
		log.Logger.Errorf("Failed to count Teams: %s", res.Error)
		return 0, false, i18n.CountModelError
	}
	return count, true, i18n.Success
}

func (t *TeamRepo) GetAll(contestID uint, limit, offset int, hidden, banned bool, preloadL ...string) ([]model.Team, int64, bool, string) {
	var (
		teams          = make([]model.Team, 0)
		count, ok, msg = t.Count(contestID, hidden, banned)
	)
	if !ok {
		return teams, count, false, msg
	}
	res := t.DB.Model(&model.Team{}).Where("contest_id = ?", contestID)
	if !hidden {
		res = res.Where("hidden = ?", false)
	}
	if !banned {
		res = res.Where("banned = ?", false)
	}
	res = preload(res, preloadL...).Limit(limit).Offset(offset).Find(&teams)
	if res.Error != nil {
		log.Logger.Errorf("Failed to get Teams: %s", res.Error)
		return teams, count, false, msg
	}
	return teams, count, true, msg
}

func (t *TeamRepo) Update(id uint, options UpdateTeamOptions) (bool, string) {
	var count int
	data := utils.UpdateOptions2Map(options)
	for {
		count++
		if count > 10 {
			log.Logger.Warningf("Failed too many times to update team due to optimistic lock")
			return false, i18n.DeadLock
		}
		team, ok, msg := t.GetByID(id)
		if !ok {
			return ok, msg
		}
		data["version"] = team.Version + 1
		res := t.DB.Model(&model.Team{}).Where("id = ? AND version = ?", id, team.Version).Updates(data)
		if res.Error != nil {
			log.Logger.Errorf("Failed to update Team: %s", res.Error)
			return false, i18n.UpdateTeamError
		}
		if res.RowsAffected == 0 {
			continue
		}
		break
	}
	return true, i18n.Success
}

func (t *TeamRepo) Delete(idL ...uint) (bool, string) {
	answerIDL, submissionIDL := make([]uint, 0), make([]uint, 0)
	for _, id := range idL {
		team, ok, msg := t.GetByID(id, "Users", "Answers", "Submissions")
		if !ok {
			return false, msg
		}
		deletedName := fmt.Sprintf("%s_deleted_%s", team.Name, utils.RandStr(6))
		if ok, msg = t.Update(id, UpdateTeamOptions{
			Name: &deletedName,
		}); !ok {
			return false, msg
		}
		for _, user := range team.Users {
			if err := DeleteUserFromContest(t.DB, user.ID, team.ContestID); err != nil {
				return false, i18n.DeleteUserFromContestError
			}
			if err := DeleteUserFromTeam(t.DB, user.ID, team.ID); err != nil {
				return false, i18n.DeleteUserFromTeamError
			}
		}
		for _, answer := range team.Answers {
			answerIDL = append(answerIDL, answer.ID)
		}
		for _, submission := range team.Submissions {
			submissionIDL = append(submissionIDL, submission.ID)
		}
	}
	if ok, msg := InitAnswerRepo(t.DB).Delete(answerIDL...); !ok {
		return false, msg
	}
	if ok, msg := InitSubmissionRepo(t.DB).Delete(submissionIDL...); !ok {
		return false, msg
	}
	if res := t.DB.Model(&model.Team{}).Where("id IN ?", idL).Delete(&model.Team{}); res.Error != nil {
		log.Logger.Errorf("Failed to delete Team: %s", res.Error)
		return false, i18n.DeleteTeamError
	}
	return true, i18n.Success
}
