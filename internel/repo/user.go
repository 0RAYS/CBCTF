package repo

import (
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
	Name     *string `json:"name"`
	Password *string `json:"password"`
	Email    *string `json:"email"`
	Country  *string `json:"country"`
	Desc     *string `json:"desc"`
	Avatar   *string `json:"avatar"`
	Verified *bool   `json:"verified"`
	Hidden   *bool   `json:"hidden"`
	Banned   *bool   `json:"banned"`
}

func InitUserRepo(tx *gorm.DB) *UserRepo {
	return &UserRepo{Repo: Repo[model.User]{DB: tx, Model: "User"}}
}

func (u *UserRepo) IsUniqueName(name string) bool {
	_, ok, _ := u.GetByName(name, false, 0)
	return !ok
}

func (u *UserRepo) IsUniqueEmail(email string) bool {
	_, ok, _ := u.GetByEmail(email, false, 0)
	return !ok
}

//func (u *UserRepo) Create(options CreateUserOptions) (model.User, bool, string) {
//	user, err := utils.S2S[model.User](options)
//	if err != nil {
//		log.Logger.Warningf("Failed to convert options to model.User: %s", err)
//		return model.User{}, false, "Options2ModelError"
//	}
//	res := u.DB.Model(&model.User{}).Create(&user)
//	if res.Error != nil {
//		log.Logger.Warningf("Failed to create User: %s", res.Error)
//		return model.User{}, false, "CreateUserError"
//	}
//	return user, true, "Success"
//}

func (u *UserRepo) getByUniqueKey(key string, value interface{}, preload bool, depth int) (model.User, bool, string) {
	switch key {
	case "name", "email":
		value = value.(string)
	case "id":
		value = value.(uint)
	default:
		return model.User{}, false, "UnsupportedKey"
	}
	var user model.User
	res := u.DB.Model(&model.User{}).Where(key+" = ?", value)
	res = model.GetPreload(res, "User", preload, depth).Find(&user).Limit(1)
	if res.RowsAffected == 0 {
		return model.User{}, false, "UserNotFound"
	}
	return user, true, "Success"
}

//func (u *UserRepo) GetByID(id uint, preload bool, depth int) (model.User, bool, string) {
//	return u.getByUniqueKey("id", id, preload, depth)
//}

func (u *UserRepo) GetByName(name string, preload bool, depth int) (model.User, bool, string) {
	return u.getByUniqueKey("name", name, preload, depth)
}

func (u *UserRepo) GetByEmail(email string, preload bool, depth int) (model.User, bool, string) {
	return u.getByUniqueKey("email", email, preload, depth)
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
		return 0, false, "CountModelError"
	}
	return count, true, "Success"
}

func (u *UserRepo) GetAll(limit, offset int, preload bool, depth int, hidden, banned bool) ([]model.User, int64, bool, string) {
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
	res = model.GetPreload(res, u.Model, preload, depth).Find(&users).Limit(limit).Offset(offset)
	if res.Error != nil {
		log.Logger.Warningf("Failed to get Users: %s", res.Error)
		return users, count, false, "GetUserError"
	}
	return users, count, true, "Success"
}

func (u *UserRepo) Update(id uint, options UpdateUserOptions) (bool, string) {
	var count int
	data := utils.UpdateOptions2Map(options)
	for {
		count++
		if count > 10 {
			log.Logger.Warningf("Failed to update User: too many times failed due to optimistic lock")
			return false, "DeadLock"
		}
		user, ok, msg := u.GetByID(id, false, 0)
		if !ok {
			return ok, msg
		}
		data["version"] = user.Version + 1
		res := u.DB.Model(&model.User{}).Omit("id", "created_at", "updated_at", "deleted_at").
			Where("id = ? AND version = ?", id, user.Version).Updates(data)
		if res.Error != nil {
			log.Logger.Warningf("Failed to update User: %s", res.Error)
			return false, "UpdateUserError"
		}
		if res.RowsAffected == 0 {
			continue
		}
		break
	}
	return true, "Success"
}

//func (u *UserRepo) Delete(idL ...uint) (bool, string) {
//	res := u.DB.Model(&model.User{}).Where("id IN ?", idL).Delete(&model.User{})
//	if res.Error != nil {
//		log.Logger.Warningf("Failed to delete User: %s", res.Error)
//		return false, "DeleteUserError"
//	}
//	return true, "Success"
//}
