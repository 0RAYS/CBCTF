package repo

import (
	"CBCTF/internel/log"
	"CBCTF/internel/model"
	"CBCTF/internel/utils"
	"gorm.io/gorm"
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
	Name    *string
	Desc    *string
	Captcha *string
	Avatar  *string
	Banned  *bool
	Hidden  *bool
}

func InitTeamRepo(tx *gorm.DB) *TeamRepo {
	return &TeamRepo{Repo: Repo[model.Team]{DB: tx, Model: "Team"}}
}

func (t *TeamRepo) IsUniqueName(contestID uint, name string) bool {
	res := t.DB.Model(&model.Team{}).Where("contest_id = ? AND name = ?", contestID, name).
		Find(&model.Team{}).Limit(1)
	return res.RowsAffected == 0
}

func (t *TeamRepo) IsUniqueMember(contestID uint, userID uint) bool {
	res := t.DB.Model(&model.UserContest{}).
		Where("contest_id = ? AND user_id = ?", contestID, userID).Find(&model.UserContest{}).Limit(1)
	return res.RowsAffected == 0
}

//func (t *TeamRepo) Create(options CreateTeamOptions) (model.Team, bool, string) {
//	team, err := utils.S2S[model.Team](options)
//	if err != nil {
//		log.Logger.Warningf("Failed to convert options to model.Team: %s", err)
//		return model.Team{}, false, "Options2ModelError"
//	}
//	if res := t.DB.Create(&team); res.Error != nil {
//		log.Logger.Errorf("Failed to create Team: %s", res.Error)
//		return model.Team{}, false, "CreateTeamError"
//	}
//	return team, true, ""
//}

func (t *TeamRepo) getByUniqueKey(key string, value interface{}, preload bool, depth int) (model.Team, bool, string) {
	switch key {
	case "id", "captain_id":
		value = value.(uint)
	default:
		return model.Team{}, false, "UnsupportedKey"
	}
	var team model.Team
	res := t.DB.Model(&model.Team{}).Where(key+" = ?", value)
	res = model.GetPreload(res, t.Model, preload, depth).Find(&team).Limit(1)
	if res.RowsAffected == 0 {
		return model.Team{}, false, "TeamNotFound"
	}
	return team, true, "Success"
}

//func (t *TeamRepo) GetByID(id uint, preload bool, depth int) (model.Team, bool, string) {
//	return t.getByUniqueKey("id", id, preload, depth)
//}

func (t *TeamRepo) GetByCaptainID(captainID uint, preload bool, depth int) (model.Team, bool, string) {
	return t.getByUniqueKey("captain_id", captainID, preload, depth)
}

// GetBy2ID 根据用户 ID 和比赛 ID 获取 model.Team, 等同于 GetByID(teamID, true, 0)
func (t *TeamRepo) GetBy2ID(userID uint, contestID uint) (model.Team, bool, string) {
	user, ok, msg := InitUserRepo(t.DB).GetByID(userID, true, 2)
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

func (t *TeamRepo) GetAll(contestID uint, limit, offset int, preload bool, depth int, hidden, banned bool) ([]model.Team, int64, bool, string) {
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
	res = model.GetPreload(res, t.Model, preload, depth).Find(&teams).Limit(limit).Offset(offset)
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
		team, ok, msg := t.GetByID(id, false, 0)
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

//func (t *TeamRepo) Delete(idL ...uint) (bool, string) {
//	res := t.DB.Model(&model.Team{}).Where("id IN ?", idL).Delete(&model.Team{})
//	if res.Error != nil {
//		log.Logger.Warningf("Failed to delete Team: %s", res.Error)
//		return false, "DeleteTeamError"
//	}
//	return true, "Success"
//}
