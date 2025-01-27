package db

import (
	"CBCTF/internal/constants"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/utils"
	"context"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// CreateUser 创建用户
func CreateUser(ctx context.Context, form constants.CreateUserForm) (model.User, bool, string) {
	if !IsValidEmail(form.Email) {
		return model.User{}, false, "InvalidEmail"
	}
	if !IsUniqueName(form.Name, model.User{}) {
		return model.User{}, false, "UserNameExists"
	}
	if !IsUniqueEmail(form.Email) {
		return model.User{}, false, "EmailExists"
	}
	user := model.InitUser(form)
	res := DB.WithContext(ctx).Model(&model.User{}).Create(&user)
	if res.Error != nil {
		log.Logger.Errorf("Failed to create user: %s", res.Error.Error())
		return model.User{}, false, "CreateUserError"
	}
	//go func() {
	//	if err := redis.DelUsersCache(); err != nil && !errors.Is(err, context.DeadlineExceeded) {
	//		log.Logger.Warningf("Failed to delete users cache: %s", err.Error())
	//	}
	//}()
	return user, true, "Success"
}

// GetUserByID 根据 ID 获取 model.User, preloadL[0] 为是否预加载, preloadL[1] 为是否嵌套预加载
func GetUserByID(ctx context.Context, id uint, preloadL ...bool) (model.User, bool, string) {
	preload := true
	nest := false
	if len(preloadL) > 0 {
		preload = preloadL[0]
	}
	if len(preloadL) > 1 {
		nest = preloadL[1]
	}
	//cacheKey := fmt.Sprintf("user:%d:%v:%v", id, preload, nest)
	//if user, ok := redis.GetUserCache(cacheKey); ok {
	//	return user, true, "Success"
	//}
	var user model.User
	res := DB.WithContext(ctx).Model(&model.User{}).Where("id = ?", id)
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
	//go func() {
	//	if err := redis.SetUserCache(cacheKey, user); err != nil && !errors.Is(err, context.DeadlineExceeded) {
	//		log.Logger.Warningf("Failed to set user cache: %s", err.Error())
	//	}
	//}()
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
	//go func() {
	//	if err := redis.DelUserCache(id); err != nil && !errors.Is(err, context.DeadlineExceeded) {
	//		log.Logger.Warningf("Failed to delete user cache: %s", err.Error())
	//	}
	//	if err := redis.DelUsersCache(); err != nil && !errors.Is(err, context.DeadlineExceeded) {
	//		log.Logger.Warningf("Failed to delete users cache: %s", err.Error())
	//	}
	//}()
	return true, "Success"
}

// UpdateUser 更新用户, 使用 map 更新属性, 结构体会导致零值未更新, 对字段值的具体要求应当交给上层实现
func UpdateUser(ctx context.Context, id uint, updateData map[string]interface{}) (bool, string) {
	res := DB.WithContext(ctx).Model(&model.User{}).Where("id = ?", id).
		Omit("id", "created_at", "updated_at", "deleted_at").Updates(updateData)
	if res.Error != nil {
		log.Logger.Errorf("Failed to update user: %s", res.Error.Error())
		return false, "UpdateUserError"
	}
	//go func() {
	//	if err := redis.DelUserCache(id); err != nil && !errors.Is(err, context.DeadlineExceeded) {
	//		log.Logger.Warningf("Failed to delete user cache: %s", err.Error())
	//	}
	//	if err := redis.DelUsersCache(); err != nil && !errors.Is(err, context.DeadlineExceeded) {
	//		log.Logger.Warningf("Failed to delete users cache: %s", err.Error())
	//	}
	//}()
	return true, "Success"
}

func VerifyUser(ctx context.Context, username string, password string) (model.User, bool, string) {
	var user model.User
	var res *gorm.DB
	res = DB.WithContext(ctx).Model(&model.User{}).Where("name = ? OR email = ?", username, username).
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

func ChangePasswordUser(ctx context.Context, id uint, oldPassword string, newPassword string) (bool, string) {
	user, ok, msg := GetUserByID(ctx, id)
	if !ok {
		return false, msg
	}
	if !utils.CompareHashAndPassword(user.Password, oldPassword) {
		return false, "PasswordError"
	}
	if utils.CompareHashAndPassword(user.Password, newPassword) {
		return false, "PasswordSame"
	}
	hash := utils.HashPassword(newPassword)
	if ok, msg = UpdateUser(ctx, id, map[string]interface{}{"password": hash}); !ok {
		return false, msg
	}
	//go func() {
	//	if err := redis.DelUserCache(id); err != nil && !errors.Is(err, context.DeadlineExceeded) {
	//		log.Logger.Warningf("Failed to delete user cache: %s", err.Error())
	//	}
	//}()
	return true, "Success"
}

func CountUsers(ctx context.Context) int64 {
	var count int64
	DB.WithContext(ctx).Model(&model.User{}).Count(&count)
	return count
}

func GetUsers(ctx context.Context, limit int, offset int, all bool, preloadL ...bool) ([]model.User, int64, bool, string) {
	if limit <= 0 {
		limit = -1
	}
	if offset <= 0 {
		offset = -1
	}
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
	res := DB.WithContext(ctx).Model(&model.User{})
	if !all {
		res = res.Where("hidden = ? AND banned = ?", false, false)
	}
	if res.Count(&count).Error != nil {
		log.Logger.Errorf("Failed to get contest count: %s", res.Error.Error())
		return nil, 0, false, "UnknownError"
	}
	//cacheKey := fmt.Sprintf("users:%v:%v:%d:%d", preload, nest, limit, offset)
	//if users, ok := redis.GetUsersCache(cacheKey); ok {
	//	return users, count, true, "Success"
	//}
	if preload {
		if nest {
			res = res.Preload("Teams.Users").Preload("Contests.Users").Preload("Contests.Teams")
		}
		res = res.Preload(clause.Associations)
	}
	if res = res.Limit(limit).Offset(offset).Find(&users); res.Error != nil {
		log.Logger.Errorf("Failed to get users: %s", res.Error.Error())
		return nil, 0, false, "UnknownError"
	}
	//go func() {
	//	if err := redis.SetUsersCache(cacheKey, users); err != nil && !errors.Is(err, context.DeadlineExceeded) {
	//		log.Logger.Warningf("Failed to set users cache: %s", err.Error())
	//	}
	//}()
	return users, count, true, "Success"

}
