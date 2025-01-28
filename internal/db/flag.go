package db

import (
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"context"
)

// GenerateFlag is a function to generate team all flags, should be call after team join contest
// 未完成
func GenerateFlag(ctx context.Context, contestID, teamID uint) {
	var (
		usages    []model.Usage
		challenge model.Challenge
		//flag      model.Flag
		ok  bool
		msg string
	)
	usages, ok, msg = GetUsageByContestID(ctx, contestID, true)
	if !ok {
		log.Logger.Warningf("Failed to get %d challenges: %s", contestID, msg)
	}
	for _, usage := range usages {
		challenge, ok, msg = GetChallengeByID(ctx, usage.ChallengeID)
		if !ok {
			log.Logger.Warningf("Failed to get challenge %s: %s", usage.ChallengeID, msg)
			continue
		}

		switch challenge.Type {
		case model.Static:
			_, ok, msg = CreateFlag(ctx, contestID, teamID, usage.ChallengeID, usage.Flag)
		case model.Dynamic:
			_, ok, msg = CreateFlag(ctx, contestID, teamID, usage.ChallengeID, "")
		case model.Container:
			_, ok, msg = CreateFlag(ctx, contestID, teamID, usage.ChallengeID, "")
		default:
			continue
		}
		if !ok {
			log.Logger.Warningf(
				"Failed to generator flag for contest %d team %d challenge %s: %s",
				contestID, teamID, usage.ChallengeID, msg,
			)
			continue
		}
	}
}

// CreateFlag is a function to create a new flag
func CreateFlag(ctx context.Context, contestID, teamID uint, challengeID, value string) (model.Flag, bool, string) {
	flag := model.InitFlag(contestID, teamID, challengeID, value)
	res := DB.WithContext(ctx).Model(model.Flag{}).Create(&flag)
	if res.Error != nil {
		log.Logger.Errorf("Failed to create Flag: %s", res.Error.Error())
		return model.Flag{}, false, "CreateFlagError"
	}
	return flag, true, "Success"
}

// GetFlagBy3ID is a function to get flag
func GetFlagBy3ID(ctx context.Context, contestID, teamID uint, challengeID string) (model.Flag, bool, string) {
	var flag model.Flag
	res := DB.WithContext(ctx).Model(model.Flag{}).
		Where("contest_id = ? AND team_id = ? AND challenge_id = ?", contestID, teamID, challengeID).Find(&flag).Limit(1)
	if res.RowsAffected != 1 {
		return model.Flag{}, false, "FlagNotFound"
	}
	return flag, true, "Success"
}

func UpdateFlag(ctx context.Context, contestID, teamID uint, challengeID, value string) (bool, string) {
	res := DB.WithContext(ctx).Model(model.Flag{}).
		Where("contest_id = ? AND team_id = ? AND challenge_id = ?", contestID, teamID, challengeID).Update("value", value)
	if res.Error != nil {
		log.Logger.Errorf("Failed to update Flag: %s", res.Error.Error())
		return false, "UpdateFlagError"
	}
	return true, "Success"
}

func VerifyFlag(ctx context.Context, contestID, teamID uint, challengeID, value string) bool {
	flag, ok, _ := GetFlagBy3ID(ctx, contestID, teamID, challengeID)
	if !ok {
		return false
	}
	if flag.Value == value {
		return true
	}
	return false
}
