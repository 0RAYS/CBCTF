package db

import (
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"context"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// CreateUser 创建用户
func CreateUser(ctx context.Context, name string, password string, email string) (model.User, bool, string) {
	if !isValidEmail(email) {
		return model.User{}, false, "InvalidEmail"
	}
	if !isUniqueName(name, model.User{}) {
		return model.User{}, false, "UserNameExists"
	}
	if !isUniqueEmail(email, model.User{}) {
		return model.User{}, false, "EmailExists"
	}
	user := model.InitUser(name, password, email)
	res := DB.WithContext(ctx).Model(&model.User{}).Create(&user)
	if res.Error != nil {
		log.Logger.Errorf("Failed to create user: %s", res.Error.Error())
		return model.User{}, false, "CreateUserError"
	}
	return user, true, "Success"
}

// GetUserByID 根据 ID 获取 model.User, preloadL[0] 为是否预加载, preloadL[1] 为是否嵌套预加载
func GetUserByID(ctx context.Context, id uint, preloadL ...bool) (model.User, bool, string) {
	var user model.User
	var res *gorm.DB
	preload := true
	nest := false
	if len(preloadL) > 0 {
		preload = preloadL[0]
	}
	if len(preloadL) > 1 {
		nest = preloadL[1]
	}
	if preload {
		if nest {
			res = DB.WithContext(ctx).Model(&model.User{}).Where("id = ?", id).
				Preload("Teams.Users").Preload("Contests.Users").Preload("Contests.Teams").
				Preload(clause.Associations).Find(&user).Limit(1)
		} else {
			res = DB.WithContext(ctx).Model(&model.User{}).Where("id = ?", id).Preload(clause.Associations).
				Find(&user).Limit(1)
		}
	} else {
		res = DB.WithContext(ctx).Model(&model.User{}).Where("id = ?", id).Find(&user).Limit(1)
	}
	if res.RowsAffected != 1 {
		return model.User{}, false, "UserNotFound"
	}
	return user, true, "Success"
}

// DeleteUser 根据 id 删除 model.User, 同时删除与 model.Team, model.Contest 的关联
func DeleteUser(ctx context.Context, id uint) (bool, string) {
	user, ok, msg := GetUserByID(ctx, id, true, true)
	if !ok {
		return false, msg
	}
	for _, team := range user.Teams {
		if len(team.Users) == 1 {
			if ok, msg = DeleteTeam(ctx, team.ID); !ok {
				log.Logger.Errorf("Failed to delete empty team: %s", msg)
				return false, msg
			}
		}
	}
	if err := DB.WithContext(ctx).Model(&model.User{}).Select(clause.Associations).Delete(&model.User{}, id).Error; err != nil {
		log.Logger.Warningf("Failed to delete user: %s", err.Error())
		return false, "DeleteUserError"
	}
	return true, "Success"
}

// UpdateUser 更新用户, 使用 map 更新属性, 结构体会导致零值未更新, 对字段值的具体要求应当交给上层实现
func UpdateUser(ctx context.Context, id uint, updateData map[string]interface{}) (bool, string) {
	res := DB.WithContext(ctx).Model(&model.User{}).Where("id = ?", id).
		Omit("id", "created_at", "updated_at", "deleted_at").Updates(updateData)
	if res.Error != nil {
		log.Logger.Errorf("Failed to update user: %s", res.Error.Error())
		return false, "UpdateError"
	}
	return true, "Success"
}
