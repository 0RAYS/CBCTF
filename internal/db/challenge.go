package db

import (
	"CBCTF/internal/constants"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"context"
	"gorm.io/gorm"
)

// CreateChallenge 创建题目
func CreateChallenge(tx *gorm.DB, form constants.CreateChallengeForm) (model.Challenge, bool, string) {
	if !IsValidChallengeType(form.Type) {
		return model.Challenge{}, false, "InvalidChallengeType"
	}
	challenge := model.InitChallenge(form)
	res := tx.Model(model.Challenge{}).Create(&challenge)
	if res.Error != nil {
		log.Logger.Warningf("Failed to create Challenge: %s", res.Error)
		return model.Challenge{}, false, "CreateChallengeError"
	}
	return challenge, true, "Success"
}

// GetChallengeByID 根据 id 获取题目
func GetChallengeByID(ctx context.Context, id string) (model.Challenge, bool, string) {
	var challenge model.Challenge
	res := DB.WithContext(ctx).Model(model.Challenge{}).Where("id = ?", id).Find(&challenge).Limit(1)
	if res.RowsAffected != 1 {
		return model.Challenge{}, false, "ChallengeNotFound"
	}
	return challenge, true, "Success"
}

// GetChallenges 获取题目列表, 可接受 type 和 category 参数
func GetChallenges(ctx context.Context, limit, offset, t int, category string) ([]model.Challenge, int64, bool, string) {
	if limit <= 0 {
		limit = -1
	}
	if offset <= 0 {
		offset = -1
	}
	var challenges []model.Challenge
	var count int64
	res := DB.WithContext(ctx).Model(model.Challenge{})
	if t != -1 && category != "" {
		res = res.Where("type = ? AND category = ?", t, category)
	} else if !(t == -1 && category == "") {
		res = res.Where("type = ? OR category = ?", t, category)
	}
	if res.Count(&count).Error != nil {
		log.Logger.Warningf("Failed to get challenge count: %v", res.Error)
		return nil, 0, false, "UnknownError"
	}
	if res = res.Limit(limit).Offset(offset).Find(&challenges); res.Error != nil {
		log.Logger.Warningf("Failed to get Challenges: %v", res.Error)
		return nil, 0, false, "UnknownError"
	}
	return challenges, count, true, "Success"
}

// CountChallenges 获取题目数量
func CountChallenges(ctx context.Context) int64 {
	var count int64
	DB.WithContext(ctx).Model(&model.Challenge{}).Count(&count)
	return count
}

// UpdateChallenge 更新题目, 使用 map 更新属性, 结构体会导致零值未更新, 对字段值的具体要求应当交给上层实现
func UpdateChallenge(tx *gorm.DB, id string, updateData map[string]interface{}) (bool, string) {
	res := tx.Model(model.Challenge{}).Where("id = ?", id).
		Omit("id", "created_at", "updated_at", "deleted_at").Updates(updateData)
	if res.Error != nil {
		log.Logger.Warningf("Failed to update Challenge: %v", res.Error)
		return false, "UpdateChallengeError"
	}
	return true, "Success"
}

// DeleteChallenge 删除题目
func DeleteChallenge(tx *gorm.DB, id string) (bool, string) {
	res := tx.Model(model.Challenge{}).Where("id = ?", id).Delete(&model.Challenge{})
	if res.Error != nil {
		log.Logger.Warningf("Failed to delete Challenge: %v", res.Error)
		return false, "DeleteChallengeError"
	}
	if !ClearByID(tx, "challenge_id", id) {
		return false, "DeleteAssociatedDataError"
	}
	return true, "Success"
}

// GetCategories 获取 type 下所有的题目分类
func GetCategories(ctx context.Context, t int) ([]string, bool, string) {
	var categories []string
	res := DB.WithContext(ctx).Model(&model.Challenge{}).Where("type = ?", t).Select("distinct category").Find(&categories)
	if res.Error != nil {
		log.Logger.Warningf("Failed to get categories: %s", res.Error)
		return nil, false, "UnknownError"
	}
	return categories, true, "Success"
}
