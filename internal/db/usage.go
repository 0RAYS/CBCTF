package db

import (
	"CBCTF/internal/constants"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"context"
)

func CreateUsage(ctx context.Context, form constants.CreateUsageForm, contestID uint) ([]model.Usage, bool, string) {
	var usages []model.Usage
	for _, c := range form.ChallengeID {
		challenge, ok, _ := GetChallengeByID(ctx, c)
		if !ok {
			log.Logger.Warningf("Failed to get challenge by ID: %s", c)
			continue
		}
		if _, ok, _ = GetUsageBy2ID(ctx, contestID, c); ok {
			continue
		}
		usage := model.InitUsage(c, contestID, challenge.Flag)
		if err := DB.WithContext(ctx).Model(model.Usage{}).Create(&usage).Error; err != nil {
			continue
		}
		usages = append(usages, usage)
	}
	return usages, true, "Success"
}

func GetUsageByContestID(ctx context.Context, contestID uint, all bool) ([]model.Usage, bool, string) {
	var usages []model.Usage
	res := DB.WithContext(ctx).Model(model.Usage{})
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

func GetUsageByChallengeID(ctx context.Context, challengeID string) ([]model.Usage, bool, string) {
	var usages []model.Usage
	res := DB.WithContext(ctx).Model(model.Usage{}).Where("challenge_id = ?", challengeID).Find(&usages)
	if res.Error != nil {
		log.Logger.Warningf("Failed to get Usage: %s", res.Error)
		return nil, false, "GetUsageError"
	}
	return usages, true, "Success"
}

func GetUsageBy2ID(ctx context.Context, contestID uint, challengeID string) (model.Usage, bool, string) {
	var usage model.Usage
	res := DB.WithContext(ctx).Model(model.Usage{}).Where("contest_id = ? AND challenge_id = ?", contestID, challengeID).Find(&usage).Limit(1)
	if res.RowsAffected != 1 {
		return model.Usage{}, false, "UsageNotFound"
	}
	return usage, true, "Success"
}

func GetUsageByID(ctx context.Context, id uint) (model.Usage, bool, string) {
	var usage model.Usage
	res := DB.WithContext(ctx).Model(model.Usage{}).Where("id = ?", id).Find(&usage).Limit(1)
	if res.RowsAffected != 1 {
		return model.Usage{}, false, "UsageNotFound"
	}
	return usage, true, "Success"
}

func UpdateUsage(ctx context.Context, id uint, updateData map[string]interface{}) (bool, string) {
	res := DB.WithContext(ctx).Model(model.Usage{}).Where("id = ?", id).
		Omit("id", "created_at", "updated_at", "deleted_at").Updates(updateData)
	if res.Error != nil {
		log.Logger.Warningf("Failed to update Usage: %s", res.Error)
		return false, "UpdateUsageError"
	}
	return true, "Success"
}

func DeleteUsage(ctx context.Context, id uint) (bool, string) {
	res := DB.WithContext(ctx).Model(model.Usage{}).Where("id = ?", id).Delete(&model.Usage{})
	if res.Error != nil {
		log.Logger.Warningf("Failed to delete Usage: %s", res.Error)
		return false, "DeleteUsageError"
	}
	return true, "Success"
}
