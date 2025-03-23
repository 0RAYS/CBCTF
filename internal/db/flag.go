package db

import (
	"CBCTF/internal/k8s"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/utils"
	"fmt"
	"gorm.io/gorm"
)

// InitFlag 生成对应 model.Team 的 model.Flag, 考虑题目类型
func InitFlag(tx *gorm.DB, contest model.Contest, team model.Team, usage model.Usage) (model.Flag, bool, string) {
	var (
		flag model.Flag
		ok   bool
		msg  string
	)
	switch usage.Type {
	case model.Static:
		flag, ok, msg = RecordFlag(tx, contest.ID, team.ID, usage.ChallengeID, fmt.Sprintf("%s{%s}", contest.Prefix, usage.Flag))
	case model.Dynamic:
		flag, ok, msg = RecordFlag(tx, contest.ID, team.ID, usage.ChallengeID, fmt.Sprintf("%s{%s}", contest.Prefix, utils.RandFlag(usage.Flag)))
		ok, msg = k8s.GenerateAttachment(usage, flag)
		if !ok {
			log.Logger.Warningf("Failed to generate flag for challenge %s: %s", flag.ChallengeID, msg)
			return model.Flag{}, false, msg
		}
	case model.Container:
		flag, ok, msg = RecordFlag(tx, contest.ID, team.ID, usage.ChallengeID, fmt.Sprintf("%s{%s}", contest.Prefix, utils.UUID()))
	default:
		flag, ok, msg = model.Flag{}, false, "InvalidChallengeType"
	}
	if !ok {
		log.Logger.Warningf(
			"Failed to generator flag for contest %d team %d challenge %s: %s",
			contest.ID, team.ID, usage.ChallengeID, msg,
		)
	}
	return flag, ok, msg
}

// RecordFlag 记录 model.Flag
func RecordFlag(tx *gorm.DB, contestID, teamID uint, challengeID, value string) (model.Flag, bool, string) {
	var (
		flag model.Flag
		ok   bool
	)
	if flag, ok, _ = GetFlagBy3ID(tx, contestID, teamID, challengeID); ok {
		ok, _ = UpdateFlag(tx, contestID, teamID, challengeID, value)
		if !ok {
			return model.Flag{}, false, "UpdateFlagError"
		}
	} else {
		flag = model.InitFlag(contestID, teamID, challengeID, value)
		res := tx.Model(&model.Flag{}).Create(&flag)
		if res.Error != nil {
			log.Logger.Warningf("Failed to create Flag: %s", res.Error)

			return model.Flag{}, false, "CreateFlagError"
		}
	}
	flag.Value = value
	return flag, true, "Success"
}

// GetFlagBy3ID 根据 contestID, teamID, challengeID 获取 model.Flag, 实际上依据 teamID 和 challengeID 即可获取唯一 model.Flag
func GetFlagBy3ID(tx *gorm.DB, contestID, teamID uint, challengeID string) (model.Flag, bool, string) {
	var flag model.Flag
	res := tx.Model(&model.Flag{}).
		Where("contest_id = ? AND team_id = ? AND challenge_id = ?", contestID, teamID, challengeID).Find(&flag).Limit(1)
	if res.RowsAffected != 1 {
		return model.Flag{}, false, "FlagNotFound"
	}
	return flag, true, "Success"
}

// UpdateFlag 更新 model.Flag 值, 固定只更新 value 字段
func UpdateFlag(tx *gorm.DB, contestID, teamID uint, challengeID, value string) (bool, string) {
	var count int
	for {
		count++
		if count > 10 {
			log.Logger.Warningf("Failed too many times to update user due to optimistic lock")
			return false, "FailedTooManyTimes"
		}
		var flag model.Flag
		res := tx.Model(&model.Flag{}).Where("contest_id = ? AND team_id = ? AND challenge_id = ?", contestID, teamID, challengeID).Find(&flag).Limit(1)
		if res.RowsAffected != 1 {
			return false, "FlagNotFound"
		}
		res = tx.Model(&flag).Update("value", value)
		if res.Error != nil {
			log.Logger.Warningf("Failed to update Flag: %s", res.Error)
			return false, "UpdateFlagError"
		}
		if res.RowsAffected == 0 {
			log.Logger.Debug("Failed to update flag due to optimistic lock")
			continue
		}
		break
	}
	return true, "Success"
}

// VerifyFlag 验证 flag 是否正确
func VerifyFlag(tx *gorm.DB, contestID, teamID uint, challengeID, value string) bool {
	flag, ok, _ := GetFlagBy3ID(tx, contestID, teamID, challengeID)
	if !ok {
		return false
	}
	if value == flag.Value {
		return true
	}
	if utils.In(value, flag.Values) {
		return true
	}
	return false
}

//func DeleteFlag(tx *gorm.DB, id uint) (bool, string) {
//	res := tx.Model(&model.Flag{}).
//		Where("id = ?", id).Delete(&model.Flag{})
//	if res.Error != nil {
//		log.Logger.Warningf("Failed to delete Flag: %s", res.Error)
//		return false, "DeleteFlagError"
//	}
//	return true, "Success"
//}
