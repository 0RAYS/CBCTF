package db

import (
	"CBCTF/internal/form"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// CreateUser 创建用户
func CreateUser(tx *gorm.DB, form form.CreateUserForm) (model.User, bool, string) {
	if !IsValidEmail(form.Email) {
		return model.User{}, false, "InvalidEmail"
	}
	if !IsUniqueName(tx, form.Name, model.User{}) {
		return model.User{}, false, "UserNameExists"
	}
	if !IsUniqueEmail(tx, form.Email) {
		return model.User{}, false, "EmailExists"
	}
	user := model.InitUser(form)
	res := tx.Model(&model.User{}).Create(&user)
	if res.Error != nil {
		log.Logger.Warningf("Failed to create user: %s", res.Error)
		return model.User{}, false, "CreateUserError"
	}
	return user, true, "Success"
}

// GetUserByID 根据 ID 获取 model.User, preloadL[0] 为是否预加载, preloadL[1] 为是否嵌套预加载
func GetUserByID(tx *gorm.DB, id uint, preloadL ...bool) (model.User, bool, string) {
	preload := true
	nest := false
	if len(preloadL) > 0 {
		preload = preloadL[0]
	}
	if len(preloadL) > 1 {
		nest = preloadL[1]
	}
	var user model.User
	res := tx.Model(&model.User{}).Where("id = ?", id)
	if preload {
		if nest {
			res = res.Preload("Teams.Users").Preload("Contests.Users").Preload("Contests.Teams")
		}
		res = res.Preload(clause.Associations)
	}
	res = res.Find(&user).Limit(1)
	if res.RowsAffected != 1 {
		return model.User{}, false, "UserNotFound"
	}
	return user, true, "Success"
}

// DeleteUser 根据 id 删除 model.User, 同时删除与 model.Team, model.Contest 的关联, 此处需嵌套预加载, 所以不接受中间件保存的值
func DeleteUser(tx *gorm.DB, id uint) (bool, string) {
	user, ok, msg := GetUserByID(tx, id, true, true)
	if !ok {
		return false, msg
	}
	for _, team := range user.Teams {
		if len(team.Users) == 1 {
			if ok, msg = DeleteTeam(tx, *team); !ok {
				log.Logger.Warningf("Failed to delete empty team: %s", msg)
				return false, msg
			}
		}
	}
	if err := tx.Model(&model.User{}).Select(clause.Associations).Delete(&model.User{}, id).Error; err != nil {
		log.Logger.Warningf("Failed to delete user: %s", err)
		return false, "DeleteUserError"
	}
	if !ClearByID(tx, "user_id", id) {
		return false, "DeleteAssociatedDataError"
	}
	return true, "Success"
}

// UpdateUser 更新用户, 使用 map 更新属性, 结构体会导致零值未更新, 对字段值的具体要求应当交给上层实现
func UpdateUser(tx *gorm.DB, id uint, updateData map[string]interface{}) (bool, string) {
	res := tx.Model(&model.User{}).Where("id = ?", id).
		Omit("id", "created_at", "updated_at", "deleted_at").Updates(updateData)
	if res.Error != nil {
		log.Logger.Warningf("Failed to update user: %s", res.Error)
		return false, "UpdateUserError"
	}
	return true, "Success"
}

// VerifyUser 验证用户
func VerifyUser(tx *gorm.DB, username string, password string) (model.User, bool, string) {
	var user model.User
	var res *gorm.DB
	res = tx.Model(&model.User{}).Where("name = ? OR email = ?", username, username).
		Find(&user).Limit(1)
	if res.RowsAffected != 1 {
		// 保持 用户名不存在 与 密码错误 行为相同
		return model.User{}, false, "NameOrPasswordError"
	}
	if utils.CompareHashAndPassword(user.Password, password) {
		return user, true, "Success"
	}
	return model.User{}, false, "NameOrPasswordError"
}

// ChangePasswordUser 修改密码
func ChangePasswordUser(tx *gorm.DB, user model.User, oldPassword string, newPassword string) (bool, string) {
	if !utils.CompareHashAndPassword(user.Password, oldPassword) {
		return false, "PasswordError"
	}
	if utils.CompareHashAndPassword(user.Password, newPassword) {
		return false, "PasswordSame"
	}
	hash := utils.HashPassword(newPassword)
	if ok, msg := UpdateUser(tx, user.ID, map[string]interface{}{"password": hash}); !ok {
		return false, msg
	}
	return true, "Success"
}

// CountUsers 获取用户数量
func CountUsers(tx *gorm.DB) int64 {
	var count int64
	tx.Model(&model.User{}).Count(&count)
	return count
}

// GetUsers 获取用户列表, 可接受 limit, offset, all 参数, preloadL[0] 为是否预加载, preloadL[1] 为是否嵌套预加载
func GetUsers(tx *gorm.DB, limit int, offset int, all bool, preloadL ...bool) ([]model.User, int64, bool, string) {
	preload := true
	nest := false
	if len(preloadL) > 0 {
		preload = preloadL[0]
	}
	if len(preloadL) > 1 {
		nest = preloadL[1]
	}
	var users []model.User
	var count int64
	res := tx.Model(&model.User{})
	if !all {
		res = res.Where("hidden = ? AND banned = ?", false, false)
	}
	if res.Count(&count).Error != nil {
		log.Logger.Warningf("Failed to get contest count: %s", res.Error)
		return nil, 0, false, "UnknownError"
	}
	if preload {
		if nest {
			res = res.Preload("Teams.Users").Preload("Contests.Users").Preload("Contests.Teams")
		}
		res = res.Preload(clause.Associations)
	}
	limit, offset = utils.TidyPaginate(int(count), limit, offset)
	if res = res.Limit(limit).Offset(offset).Find(&users); res.Error != nil {
		log.Logger.Warningf("Failed to get users: %s", res.Error)
		return nil, 0, false, "UnknownError"
	}
	return users, count, true, "Success"

}
