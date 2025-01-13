package db

import (
	"RayWar/internal/log"
	"RayWar/internal/model"
	"RayWar/internal/utils"
	"gorm.io/gorm"
)

// CreateUser 创建新用户，单独判断用户名邮箱是否合法并且唯一
func CreateUser(name string, password string, email string) (model.User, bool, string) {
	if !isValidEmail(email) {
		return model.User{}, false, "InvalidEmail"
	}
	if !isUniqueUserName(name) {
		return model.User{}, false, "UserNameExists"
	}
	if !isUniqueEmail(email) {
		return model.User{}, false, "EmailExists"
	}
	user := model.InitUser(name, password, email)
	res := DB.Model(&model.User{}).Create(&user)
	if res.Error != nil {
		log.Logger.Warningf("Failed to create user: %s", res.Error.Error())
		return model.User{}, false, "CreateUserError"
	}
	return user, true, "Success"
}

// GetUserByName 根据 Name 获取 model.User
func GetUserByName(name string, preloadL ...bool) (model.User, bool, string) {
	var user model.User
	var res *gorm.DB
	preload := true
	if len(preloadL) > 0 {
		preload = preloadL[0]
	}
	if preload {
		res = DB.Model(&model.User{}).Where("name = ?", name).Preload("Teams").
			Find(&user).Limit(1)
	} else {
		res = DB.Model(&model.User{}).Where("name = ?", name).Find(&user).Limit(1)
	}
	if res.RowsAffected != 1 {
		return model.User{}, false, "UserNotFound"
	}
	return user, true, "Success"
}

// GetUserByID 根据 ID 获取 model.User
func GetUserByID(id uint, preloadL ...bool) (model.User, bool, string) {
	var user model.User
	var res *gorm.DB
	preload := true
	if len(preloadL) > 0 {
		preload = preloadL[0]
	}
	if preload {
		res = DB.Model(&model.User{}).Where("id = ?", id).Preload("Teams").
			Find(&user).Limit(1)
	} else {
		res = DB.Model(&model.User{}).Where("id = ?", id).Find(&user).Limit(1)
	}
	if res.RowsAffected != 1 {
		return model.User{}, false, "UserNotFound"
	}
	return user, true, "Success"
}

// VerifyUser 根据用户名或邮箱和密码获取用户，主要用于登录
func VerifyUser(username string, password string) (model.User, bool, string) {
	var user model.User
	var res *gorm.DB
	res = DB.Model(&model.User{}).Where("name = ? OR email = ?", username, username).
		Preload("Teams").Find(&user).Limit(1)
	if res.RowsAffected != 1 {
		// 保持 用户名不存在 与 密码错误 行为相同
		return model.User{}, false, "NameOrPasswordError"
	}
	if utils.CompareHashAndPassword(user.Password, password) {
		return user, true, "Success"
	}
	return model.User{}, false, "NameOrPasswordError"
}

// UpdateUser 对字段值的具体要求应当交给上层实现
func UpdateUser(user model.User, updateData map[string]interface{}) (bool, string) {
	res := DB.Model(&model.User{}).Where("id = ?", user.ID).Omit("id", "created_at", "updated_at", "deleted_at").Updates(updateData)
	if res.Error != nil {
		log.Logger.Warningf("Failed to create user: %s", res.Error.Error())
		return false, "UpdateError"
	}
	return true, "Success"
}

// DeleteUser 同时删除与 model.Team 的关联关系，但不删除关联的数据，可能导致某个队伍人数为0，定义 ClearEmptyTeam
func DeleteUser(user model.User) (bool, string) {
	if err := DB.Model(&model.User{}).Select("Teams").Delete(&user).Error; err != nil {
		log.Logger.Warningf("Failed to delete user: %s", err.Error())
		return false, "DeleteUserError"
	}
	return true, "Success"
}

// JoinTeam model.User 加入 model.Team，先鉴别是否是三个值是否合法
func JoinTeam(user model.User, contest model.Contest, team model.Team) (bool, string) {
	if !isNotRepeatPlayer(contest.ID, user.ID) {
		return false, "RepeatPlayer"
	}
	if !isContestTeam(team, contest) {
		return false, "TeamNotFound"
	}
	// 用户 团队 被ban时不能加入，但是hidden可以，比赛hidden时不可以加入
	if user.Banned || user.Type == "admin" || team.Banned || contest.Hidden {
		return false, "Forbidden"
	}
	if len(team.Users)+1 > int(contest.Size) {
		return false, "TeamIsFull"
	}
	err := DB.Model(&team).Association("Users").Append(&user)
	if err != nil {
		log.Logger.Warningf("Failed to join team: %s", err.Error())
		return false, "JoinTeamError"
	}
	return true, "Success"
}

// LeaveTeam model.User 离开 model.Team，先鉴别是否是三个值是否合法
func LeaveTeam(user model.User, contest model.Contest, team model.Team) (bool, string) {
	if isNotRepeatPlayer(contest.ID, user.ID) {
		return false, "IsNotPlayer"
	}
	if !isContestTeam(team, contest) || !isTeamUser(user, team) {
		return false, "TeamNotFound"
	}
	err := DB.Model(&team).Association("Users").Delete(user)
	if err != nil {
		log.Logger.Warningf("Failed to leave team: %s", err.Error())
		return false, "LeaveTeamError"
	}
	return true, "Success"
}

// GetUsers 获取所有用户
func GetUsers(limit int, offset int, all bool) ([]model.User, int64, bool, string) {
	if limit <= 0 {
		limit = -1
	}
	if offset <= 0 {
		offset = -1
	}
	var users []model.User
	var total int64
	if all {
		if res := DB.Model(&model.User{}).Count(&total); res.Error != nil {
			log.Logger.Warningf("Failed to get users: %s", res.Error.Error())
			return nil, 0, false, "UnknownError"
		}
		if res := DB.Model(&model.User{}).Limit(limit).Offset(offset).Find(&users); res.Error != nil {
			log.Logger.Warningf("Failed to get users: %s", res.Error.Error())
			return nil, 0, false, "UserNotFound"
		}
	} else {
		if res := DB.Model(&model.User{}).Where("hidden != ? and banned != ?", true, true).
			Count(&total); res.Error != nil {
			log.Logger.Warningf("Failed to get users: %s", res.Error.Error())
			return nil, 0, false, "UnknownError"
		}
		if res := DB.Model(&model.User{}).Where("hidden != ? and banned != ?", true, true).
			Limit(limit).Offset(offset).Find(&users); res.Error != nil {
			log.Logger.Warningf("Failed to get users: %s", res.Error.Error())
			return nil, 0, false, "UserNotFound"
		}
	}
	return users, total, true, "Success"
}
