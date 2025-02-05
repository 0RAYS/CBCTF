package db

import (
	"CBCTF/internal/k8s"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/utils"
	"fmt"
	"gorm.io/gorm"
)

// InitFlag is a function to generate team flag
func InitFlag(tx *gorm.DB, contest model.Contest, team model.Team, usage model.Usage) (model.Flag, bool, string) {
	var (
		challenge model.Challenge
		flag      model.Flag
		ok        bool
		msg       string
	)
	challenge, ok, msg = GetChallengeByID(tx, usage.ChallengeID)
	if !ok {
		return model.Flag{}, false, msg
	}
	switch challenge.Type {
	case model.Static:
		flag, ok, msg = RecordFlag(tx, contest.ID, team.ID, usage.ChallengeID, fmt.Sprintf("%s{%s}", contest.Prefix, challenge.Flag))
	case model.Dynamic:
		flag, ok, msg = RecordFlag(tx, contest.ID, team.ID, usage.ChallengeID, fmt.Sprintf("%s{%s}", contest.Prefix, utils.RandFlag(challenge.Flag)))
		go func(c model.Challenge, f model.Flag) {
			log.Logger.Debugf("Generating attachment for team %d challenge %s", team.ID, usage.ChallengeID)
			ok, msg = k8s.GenerateAttachment(c, f)
			if !ok {
				log.Logger.Warningf("Failed to generate flag for challenge %s: %s", usage.ChallengeID, msg)
			}
		}(challenge, flag)
	case model.Container:
		flag, ok, msg = RecordFlag(tx, contest.ID, team.ID, usage.ChallengeID, fmt.Sprintf("%s{%s}", contest.Prefix, utils.RandomString()))
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

// RecordFlag is a function to create a new flag
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
		res := tx.Model(model.Flag{}).Create(&flag)
		if res.Error != nil {
			log.Logger.Warningf("Failed to create Flag: %s", res.Error)

			return model.Flag{}, false, "CreateFlagError"
		}
	}
	flag.Value = value
	return flag, true, "Success"
}

// GetFlagBy3ID is a function to get flag
func GetFlagBy3ID(tx *gorm.DB, contestID, teamID uint, challengeID string) (model.Flag, bool, string) {
	var flag model.Flag
	res := tx.Model(model.Flag{}).
		Where("contest_id = ? AND team_id = ? AND challenge_id = ?", contestID, teamID, challengeID).Find(&flag).Limit(1)
	if res.RowsAffected != 1 {
		return model.Flag{}, false, "FlagNotFound"
	}
	return flag, true, "Success"
}

// UpdateFlag is a function to update flag
func UpdateFlag(tx *gorm.DB, contestID, teamID uint, challengeID, value string) (bool, string) {
	res := tx.Model(model.Flag{}).
		Where("contest_id = ? AND team_id = ? AND challenge_id = ?", contestID, teamID, challengeID).Update("value", value)
	if res.Error != nil {
		log.Logger.Warningf("Failed to update Flag: %s", res.Error)
		return false, "UpdateFlagError"
	}
	return true, "Success"
}

// VerifyFlag is a function to verify flag
func VerifyFlag(tx *gorm.DB, contestID, teamID uint, challengeID, value string) bool {
	flag, ok, _ := GetFlagBy3ID(tx, contestID, teamID, challengeID)
	if !ok {
		return false
	}
	if flag.Value == value {
		return true
	}
	return false
}
