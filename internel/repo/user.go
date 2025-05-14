package repo

import (
	"CBCTF/internel/i18n"
	"CBCTF/internel/log"
	"CBCTF/internel/model"
	"CBCTF/internel/utils"
	"gorm.io/gorm"
)

type UserRepo struct {
	Repo[model.User]
}

type CreateUserOptions struct {
	Name     string
	Password string
	Email    string
	Country  string
	Avatar   string
	Desc     string
	Verified bool
	Hidden   bool
	Banned   bool
}

type UpdateUserOptions struct {
	Name     *string  `json:"name"`
	Password *string  `json:"password"`
	Email    *string  `json:"email"`
	Country  *string  `json:"country"`
	Desc     *string  `json:"desc"`
	Avatar   *string  `json:"avatar"`
	Verified *bool    `json:"verified"`
	Hidden   *bool    `json:"hidden"`
	Banned   *bool    `json:"banned"`
	Score    *float64 `json:"score"`
	Solved   *int64   `json:"solved"`
}

func InitUserRepo(tx *gorm.DB) *UserRepo {
	return &UserRepo{Repo: Repo[model.User]{DB: tx, Model: "User"}}
}

func (u *UserRepo) IsUniqueName(name string) bool {
	_, ok, _ := u.GetByName(name)
	return !ok
}

func (u *UserRepo) IsUniqueEmail(email string) bool {
	_, ok, _ := u.GetByEmail(email)
	return !ok
}

func (u *UserRepo) getByUniqueKey(key string, value interface{}, preloadL ...string) (model.User, bool, string) {
	switch key {
	case "name", "email":
		value = value.(string)
	case "id":
		value = value.(uint)
	default:
		return model.User{}, false, i18n.UnsupportedKey
	}
	var user model.User
	res := u.DB.Model(&model.User{}).Where(key+" = ?", value)
	res = preload(res, preloadL...).Limit(1).Find(&user)
	if res.Error != nil {
		log.Logger.Warningf("Failed to get User")
		return model.User{}, false, i18n.GetUserError
	}
	if res.RowsAffected == 0 {
		return model.User{}, false, i18n.UserNotFound
	}
	return user, true, i18n.Success
}

func (u *UserRepo) GetByName(name string, preloadL ...string) (model.User, bool, string) {
	return u.getByUniqueKey("name", name, preloadL...)
}

func (u *UserRepo) GetByEmail(email string, preloadL ...string) (model.User, bool, string) {
	return u.getByUniqueKey("email", email, preloadL...)
}

func (u *UserRepo) Count(hidden, banned bool) (int64, bool, string) {
	var count int64
	res := u.DB.Model(&model.User{})
	if !hidden {
		res = res.Where("hidden = ?", false)
	}
	if !banned {
		res = res.Where("banned = ?", false)
	}
	res = res.Count(&count)
	if res.Error != nil {
		log.Logger.Warningf("Failed to count Users: %s", res.Error)
		return 0, false, i18n.CountModelError
	}
	return count, true, i18n.Success
}

func (u *UserRepo) GetAll(limit, offset int, hidden, banned bool, preloadL ...string) ([]model.User, int64, bool, string) {
	var (
		users          = make([]model.User, 0)
		count, ok, msg = u.Count(hidden, banned)
	)
	if !ok {
		return users, count, false, msg
	}
	res := u.DB.Model(&model.User{})
	if !hidden {
		res = res.Where("banned = ?", false)
	}
	if !banned {
		res = res.Where("hidden = ?", false)
	}
	res = preload(res, preloadL...).Limit(limit).Offset(offset).Find(&users)
	if res.Error != nil {
		log.Logger.Warningf("Failed to get Users: %s", res.Error)
		return users, count, false, i18n.GetUserError
	}
	return users, count, true, i18n.Success
}

func (u *UserRepo) Update(id uint, options UpdateUserOptions) (bool, string) {
	var count int
	data := utils.UpdateOptions2Map(options)
	for {
		count++
		if count > 10 {
			log.Logger.Warningf("Failed to update User: too many times failed due to optimistic lock")
			return false, i18n.DeadLock
		}
		user, ok, msg := u.GetByID(id)
		if !ok {
			return ok, msg
		}
		data["version"] = user.Version + 1
		res := u.DB.Model(&model.User{}).Where("id = ? AND version = ?", id, user.Version).Updates(data)
		if res.Error != nil {
			log.Logger.Warningf("Failed to update User: %s", res.Error)
			return false, i18n.UpdateUserError
		}
		if res.RowsAffected == 0 {
			continue
		}
		break
	}
	return true, i18n.Success
}

func (u *UserRepo) Delete(idL ...uint) (bool, string) {
	submissionIDL := make([]uint, 0)
	for _, id := range idL {
		user, ok, msg := u.GetByID(id, "Teams", "Submissions")
		if !ok {
			return false, msg
		}
		for _, team := range user.Teams {
			if err := DeleteUserFromContest(u.DB, user.ID, team.ContestID); err != nil {
				return false, i18n.DeleteUserFromContestError
			}
			if err := DeleteUserFromTeam(u.DB, user.ID, team.ID); err != nil {
				return false, i18n.DeleteUserFromTeamError
			}
		}
		for _, submission := range user.Submissions {
			submissionIDL = append(submissionIDL, submission.ID)
		}
	}
	if ok, msg := InitSubmissionRepo(u.DB).Delete(submissionIDL...); !ok {
		return false, msg
	}
	if res := u.DB.Model(&model.User{}).Where("id IN ?", idL).Delete(&model.Submission{}); res.Error != nil {
		log.Logger.Warningf("Failed to delete User: %s", res.Error)
		return false, i18n.DeleteUserError
	}
	return true, i18n.Success
}
