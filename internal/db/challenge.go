package db

import (
	"CBCTF/internal/config"
	"CBCTF/internal/constants"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/utils"
	"context"
	"fmt"
)

func CreateChallenge(ctx context.Context, form constants.CreateChallengeForm) (model.Challenge, bool, string) {
	if !IsValidChallengeType(form.Type) {
		return model.Challenge{}, false, "InvalidChallengeType"
	}
	path := fmt.Sprintf("%s/challenges/%s", config.Env.Gin.Upload.Path, utils.RandomString())
	challenge := model.InitChallenge(form, path)
	result := DB.WithContext(ctx).Model(model.Challenge{}).Create(&challenge)
	if result.Error != nil {
		log.Logger.Errorf("Failed to create Challenge: %s", result.Error.Error())
		return model.Challenge{}, false, "CreateChallengeError"
	}
	return challenge, true, "Success"
}

func GetChallengeByID(ctx context.Context, id string) (model.Challenge, bool, string) {
	var challenge model.Challenge
	result := DB.WithContext(ctx).Model(model.Challenge{}).Where("id = ?", id).Find(&challenge).Limit(1)
	if result.RowsAffected != 1 {
		return model.Challenge{}, false, "ChallengeNotFound"
	}
	return challenge, true, "Success"
}

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
	if t != 0 && category != "" {
		res = res.Where("type = ? AND category = ?", t, category)
	} else if !(t == 0 && category == "") {
		res = res.Where("type = ? OR category = ?", t, category)
	}
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

func CountChallenges(ctx context.Context) int64 {
	var count int64
	DB.WithContext(ctx).Model(&model.Challenge{}).Count(&count)
	return count
}

func UpdateChallenge(ctx context.Context, id string, updateData map[string]interface{}) (bool, string) {
	result := DB.WithContext(ctx).Model(model.Challenge{}).Where("id = ?", id).
		Omit("id", "created_at", "updated_at", "deleted_at").Updates(updateData)
	if result.Error != nil {
		log.Logger.Warningf("Failed to update Challenge: %v", result.Error.Error())
		return false, "UpdateChallengeError"
	}
	return true, "Success"
}

func DeleteChallenge(ctx context.Context, id string) (bool, string) {
	result := DB.WithContext(ctx).Model(model.Challenge{}).Where("id = ?", id).Delete(&model.Challenge{})
	if result.Error != nil {
		log.Logger.Warningf("Failed to delete Challenge: %v", result.Error.Error())
		return false, "DeleteChallengeError"
	}
	return true, "Success"
}

func GetCategories(ctx context.Context, t int) ([]string, bool, string) {
	var categories []string
	res := DB.WithContext(ctx).Model(&model.Challenge{}).Where("type = ?", t).Select("distinct category").Find(&categories)
	if res.Error != nil {
		log.Logger.Errorf("Failed to get categories: %s", res.Error.Error())
		return nil, false, "UnknownError"
	}
	return categories, true, "Success"
}
