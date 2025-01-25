package db

import (
	"CBCTF/internal/constants"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"context"
)

func CreateChallenge(ctx context.Context, form constants.CreateChallengeForm) (model.Challenge, bool, string) {
	if !IsValidChallengeType(form.Type) {
		return model.Challenge{}, false, "InvalidChallengeType"
	}
	challenge := model.InitChallenge(form)
	result := DB.WithContext(ctx).Model(model.Challenge{}).Create(&challenge)
	if result.Error != nil {
		log.Logger.Errorf("Failed to create Challenge: %s", result.Error.Error())
		return model.Challenge{}, false, "CreateChallengeError"
	}
	return challenge, true, "Success"
}

func GetChallengeByID(ctx context.Context, id uint) (model.Challenge, bool, string) {
	var challenge model.Challenge
	result := DB.WithContext(ctx).Model(model.Challenge{}).Where("id = ?", id).Find(&challenge)
	if result.RowsAffected != 1 {
		return model.Challenge{}, false, "ChallengeNotFound"
	}
	return challenge, true, "Success"
}

func GetChallenges(ctx context.Context, limit, offset int) ([]model.Challenge, int64, bool, string) {
	var challenges []model.Challenge
	var count int64
	res := DB.WithContext(ctx).Model(model.Challenge{})
	if res.Count(&count).Error != nil {
		log.Logger.Warningf("Failed to get challenge count: %v", res.Error.Error())
		return nil, 0, false, "UnknownError"
	}
	if res = res.Limit(limit).Offset(offset).Find(&challenges); res.Error != nil {
		log.Logger.Warningf("Failed to get Challenges: %v", res.Error.Error())
		return nil, 0, false, "UnknownError"
	}
	return challenges, count, true, "Success"
}

func UpdateChallenge(ctx context.Context, id uint, updateData map[string]interface{}) (bool, string) {
	result := DB.WithContext(ctx).Model(model.Challenge{}).Where("id = ?", id).
		Omit("id", "created_at", "updated_at", "deleted_at").Updates(updateData)
	if result.Error != nil {
		log.Logger.Warningf("Failed to update Challenge: %v", result.Error.Error())
		return false, "UpdateChallengeError"
	}
	return true, "Success"
}

func DeleteChallenge(ctx context.Context, id uint) (bool, string) {
	result := DB.WithContext(ctx).Model(model.Challenge{}).Where("id = ?", id).Delete(&model.Challenge{})
	if result.Error != nil {
		log.Logger.Warningf("Failed to delete Challenge: %v", result.Error.Error())
		return false, "DeleteChallengeError"
	}
	return true, "Success"
}
