package repo

import (
	"CBCTF/internel/log"
	"CBCTF/internel/model"
	"CBCTF/internel/utils"
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

func (t *TeamRepo) getByUniqueKey(key string, value interface{}, preloadL ...string) (model.Team, bool, string) {
	switch key {
	case "id", "captain_id":
		value = value.(uint)
	default:
		return model.Team{}, false, "UnsupportedKey"
	}
	var team model.Team
	res := t.DB.Model(&model.Team{}).Where(key+" = ?", value)
	res = preload(res, preloadL...).Limit(1).Find(&team)
	if res.RowsAffected == 0 {
		return model.Team{}, false, "TeamNotFound"
	}
	return team, true, "Success"
}

func (t *TeamRepo) GetByCaptainID(captainID uint, preloadL ...string) (model.Team, bool, string) {
	return t.getByUniqueKey("captain_id", captainID, preloadL...)
}

func (t *TeamRepo) GetByName(contestID uint, name string, preloadL ...string) (model.Team, bool, string) {
	var team model.Team
	res := t.DB.Model(&model.Team{}).Where("contest_id = ? AND name = ?", contestID, name)
	res = preload(res, preloadL...).Limit(1).Find(&team)
	if res.RowsAffected == 0 {
		return model.Team{}, false, "TeamNotFound"
	}
	return team, true, "Success"
}

// GetBy2ID 根据用户 ID 和比赛 ID 获取 model.Team, 等同于 GetByID(teamID, true, 0)
func (t *TeamRepo) GetBy2ID(userID uint, contestID uint) (model.Team, bool, string) {
	user, ok, msg := InitUserRepo(t.DB).
		GetByID(
			userID,
			"all",
			"Teams.Contest", "Teams.Users", "Teams.Answers", "Teams.Submissions",
			"Teams.Containers", "Teams.Cheats",
		)
	if !ok {
		return model.Team{}, false, msg
	}
	for _, team := range user.Teams {
		if team.ContestID == contestID {
			return *team, true, "Success"
		}
	}
	return model.Team{}, false, "UserNotInTeam"
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
		return 0, false, "CountModelError"
	}
	return count, true, "Success"
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
			return false, "DeadLock"
		}
		team, ok, msg := t.GetByID(id)
		if !ok {
			return ok, msg
		}
		data["version"] = team.Version + 1
		res := t.DB.Model(&model.Team{}).Omit("id", "created_at", "updated_at", "deleted_at").
			Where("id = ? AND version = ?", id, team.Version).Updates(data)
		if res.Error != nil {
			log.Logger.Errorf("Failed to update Team: %s", res.Error)
			return false, "UpdateTeamError"
		}
		if res.RowsAffected == 0 {
			continue
		}
		break
	}
	return true, "Success"
}
