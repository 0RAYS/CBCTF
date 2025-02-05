package db

import (
	"CBCTF/internal/constants"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"gorm.io/gorm"
)

// CreateUsage 创建将题目添加至比赛的记录
func CreateUsage(tx *gorm.DB, form constants.CreateUsageForm, contestID uint) ([]model.Usage, bool, string) {
	var usages []model.Usage
	for _, c := range form.ChallengeID {
		challenge, ok, _ := GetChallengeByID(tx, c)
		if !ok {
			log.Logger.Warningf("Failed to get challenge by ID: %s", c)
			continue
		}
		if _, ok, _ = GetUsageBy2ID(tx, contestID, c); ok {
			continue
		}
		usage := model.InitUsage(c, contestID, challenge.Flag)
		// 如果创建失败则跳过, 不回滚
		if err := tx.Model(model.Usage{}).Create(&usage).Error; err != nil {
			continue
		}
		usages = append(usages, usage)
	}
	return usages, true, "Success"
}

// GetUsageByContestID 获取引用
func GetUsageByContestID(tx *gorm.DB, contestID uint, all bool) ([]model.Usage, bool, string) {
	var usages []model.Usage
	res := tx.Model(model.Usage{})
	if all {
		res = res.Where("contest_id = ?", contestID)
	} else {
		res = res.Where("contest_id = ? AND hidden = ?", contestID, false)
	}
	if res := res.Find(&usages); res.Error != nil {
		log.Logger.Warningf("Failed to get Usage: %s", res.Error)
		return nil, false, "GetUsageError"
	}
	return usages, true, "Success"
}

// GetUsageByChallengeID 获取引用
func GetUsageByChallengeID(tx *gorm.DB, challengeID string) ([]model.Usage, bool, string) {
	var usages []model.Usage
	res := tx.Model(model.Usage{}).Where("challenge_id = ?", challengeID).Find(&usages)
	if res.Error != nil {
		log.Logger.Warningf("Failed to get Usage: %s", res.Error)
		return nil, false, "GetUsageError"
	}
	return usages, true, "Success"
}

// GetUsageBy2ID 获取引用
func GetUsageBy2ID(tx *gorm.DB, contestID uint, challengeID string) (model.Usage, bool, string) {
	var usage model.Usage
	res := tx.Model(model.Usage{}).Where("contest_id = ? AND challenge_id = ?", contestID, challengeID).Find(&usage).Limit(1)
	if res.RowsAffected != 1 {
		return model.Usage{}, false, "UsageNotFound"
	}
	return usage, true, "Success"
}

// GetUsageByID 获取引用
func GetUsageByID(tx *gorm.DB, id uint) (model.Usage, bool, string) {
	var usage model.Usage
	res := tx.Model(model.Usage{}).Where("id = ?", id).Find(&usage).Limit(1)
	if res.RowsAffected != 1 {
		return model.Usage{}, false, "UsageNotFound"
	}
	return usage, true, "Success"
}

// UpdateUsage 更新引用
func UpdateUsage(tx *gorm.DB, id uint, updateData map[string]interface{}) (bool, string) {
	res := tx.Model(model.Usage{}).Where("id = ?", id).
		Omit("id", "created_at", "updated_at", "deleted_at").Updates(updateData)
	if res.Error != nil {
		log.Logger.Warningf("Failed to update Usage: %s", res.Error)
		return false, "UpdateUsageError"
	}
	return true, "Success"
}

func AddSolvers(tx *gorm.DB, id uint) (bool, string) {
	res := tx.Model(model.Usage{}).Where("id = ?", id).
		UpdateColumn("solvers", gorm.Expr("solvers + ?", 1))
	if res.Error != nil {
		log.Logger.Warningf("Failed to update Usage: %s", res.Error)
		return false, "UpdateUsageError"
	}
	return true, "Success"
}

// DeleteUsage 删除引用
func DeleteUsage(tx *gorm.DB, id uint) (bool, string) {
	res := tx.Model(model.Usage{}).Where("id = ?", id).Delete(&model.Usage{})
	if res.Error != nil {
		log.Logger.Warningf("Failed to delete Usage: %s", res.Error)
		return false, "DeleteUsageError"
	}
	return true, "Success"
}
