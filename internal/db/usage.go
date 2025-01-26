package db

import (
	"CBCTF/internal/constants"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"context"
)

func CreateUsage(ctx context.Context, form constants.CreateUsageForm, contestID uint) (model.Usage, bool, string) {
	challenge, ok, msg := GetChallengeByID(ctx, form.ChallengeID)
	if !ok {
		return model.Usage{}, false, msg
	}
	usage := model.InitUsage(form, contestID, challenge.Flag)
	if err := DB.WithContext(ctx).Model(model.Usage{}).Create(&usage).Error; err != nil {
		log.Logger.Errorf("Failed to create Usage: %s", err.Error())
		return model.Usage{}, false, "CreateUsageError"
	}
	return usage, true, "Success"
}

func GetUsageByContestID(ctx context.Context, contestID uint) ([]model.Usage, bool, string) {
	var usages []model.Usage
	res := DB.WithContext(ctx).Model(model.Usage{}).Where("contest_id = ?", contestID).Find(&usages)
	if res.Error != nil {
		log.Logger.Errorf("Failed to get Usage: %s", res.Error.Error())
		return nil, false, "GetUsageError"
	}
	return usages, true, "Success"
}

func GetUsageByChallengeID(ctx context.Context, challengeID string) ([]model.Usage, bool, string) {
	var usages []model.Usage
	res := DB.WithContext(ctx).Model(model.Usage{}).Where("challenge_id = ?", challengeID).Find(&usages)
	if res.Error != nil {
		log.Logger.Errorf("Failed to get Usage: %s", res.Error.Error())
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
		log.Logger.Errorf("Failed to update Usage: %s", res.Error.Error())
		return false, "UpdateUsageError"
	}
	return true, "Success"
}

func DeleteUsage(ctx context.Context, id uint) (bool, string) {
	res := DB.WithContext(ctx).Model(model.Usage{}).Where("id = ?", id).Delete(&model.Usage{})
	if res.Error != nil {
		log.Logger.Errorf("Failed to delete Usage: %s", res.Error.Error())
		return false, "DeleteUsageError"
	}
	return true, "Success"
}
