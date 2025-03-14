package db

import (
	"CBCTF/internal/form"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/utils"
	"gorm.io/gorm"
)

// CreateChallenge 创建题目
func CreateChallenge(tx *gorm.DB, form form.CreateChallengeForm) (model.Challenge, bool, string) {
	if !IsValidChallengeType(form.Type) {
		return model.Challenge{}, false, "InvalidChallengeType"
	}
	challenge := model.InitChallenge(form)
	res := tx.Model(&model.Challenge{}).Create(&challenge)
	if res.Error != nil {
		log.Logger.Warningf("Failed to create Challenge: %s", res.Error)
		return model.Challenge{}, false, "CreateChallengeError"
	}
	return challenge, true, "Success"
}

// GetChallengeByID 根据 id 获取题目
func GetChallengeByID(tx *gorm.DB, id string) (model.Challenge, bool, string) {
	var challenge model.Challenge
	res := tx.Model(&model.Challenge{}).Where("id = ?", id).Find(&challenge).Limit(1)
	if res.RowsAffected != 1 {
		return model.Challenge{}, false, "ChallengeNotFound"
	}
	return challenge, true, "Success"
}

// GetChallenges 获取题目列表, 可接受 type 和 category 参数
func GetChallenges(tx *gorm.DB, limit, offset int, t string, category string) ([]model.Challenge, int64, bool, string) {
	var challenges []model.Challenge
	var count int64
	res := tx.Model(&model.Challenge{})
	if t != "" && category != "" {
		res = res.Where("type = ? AND category = ?", t, category)
	} else if !(t == "" && category == "") {
		res = res.Where("type = ? OR category = ?", t, category)
	}
	if res.Count(&count).Error != nil {
		log.Logger.Warningf("Failed to get challenge count: %v", res.Error)
		return make([]model.Challenge, 0), 0, false, "UnknownError"
	}
	limit, offset = utils.TidyPaginate(int(count), limit, offset)
	if res = res.Limit(limit).Offset(offset).Find(&challenges); res.Error != nil {
		log.Logger.Warningf("Failed to get Challenges: %v", res.Error)
		return make([]model.Challenge, 0), 0, false, "UnknownError"
	}
	return challenges, count, true, "Success"
}

// CountChallenges 获取题目数量
func CountChallenges(tx *gorm.DB) int64 {
	var count int64
	tx.Model(&model.Challenge{}).Count(&count)
	return count
}

// UpdateChallenge 更新题目, 使用 map 更新属性, 结构体会导致零值未更新, 对字段值的具体要求应当交给上层实现
func UpdateChallenge(tx *gorm.DB, id string, updateData map[string]interface{}) (bool, string) {
	var count int
	for {
		count++
		if count > 10 {
			log.Logger.Warningf("Failed too many times to update user due to optimistic lock")
			return false, "FailedTooManyTimes"
		}
		var challenge model.Challenge
		res := tx.Model(&model.Challenge{}).Where("id = ?", id).Find(&challenge).Limit(1)
		if res.RowsAffected != 1 {
			return false, "ChallengeNotFound"
		}
		res = tx.Model(&challenge).Omit("id", "created_at", "updated_at", "deleted_at").Updates(updateData)
		if res.Error != nil {
			log.Logger.Warningf("Failed to update Challenge: %v", res.Error)
			return false, "UpdateChallengeError"
		}
		if res.RowsAffected == 0 {
			log.Logger.Debug("Failed to update challenge due to optimistic lock")
			continue
		}
		break
	}
	return true, "Success"
}

// DeleteChallenge 删除题目
func DeleteChallenge(tx *gorm.DB, id string) (bool, string) {
	res := tx.Model(&model.Challenge{}).Where("id = ?", id).Delete(&model.Challenge{})
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
func GetCategories(tx *gorm.DB, t string) ([]string, bool, string) {
	var categories []string
	res := tx.Model(&model.Challenge{})
	if t != "" {
		res = res.Where("type = ?", t)
	}
	res = res.Select("distinct category").Find(&categories)
	if res.Error != nil {
		log.Logger.Warningf("Failed to get categories: %s", res.Error)
		return make([]string, 0), false, "UnknownError"
	}
	return categories, true, "Success"
}
