package service

import (
	"CBCTF/internal/db"
	"CBCTF/internal/model"

	"gorm.io/gorm"
)

func AttachmentPath(tx *gorm.DB, challenge model.Challenge, teamID uint) (string, model.RetVal) {
	if challenge.Type != model.DynamicChallengeType {
		return challenge.AttachmentPath(teamID), model.SuccessRetVal()
	}
	var values []string
	if teamID == 0 {
		flags, _, ret := db.InitChallengeFlagRepo(tx).List(-1, -1, db.GetOptions{Conditions: map[string]any{"challenge_id": challenge.ID}, Sort: []string{"id ASC"}})
		if !ret.OK {
			return "", ret
		}
		for _, flag := range flags {
			values = append(values, flag.Value)
		}
	} else {
		ids := tx.Model(&model.ChallengeFlag{}).Select("id").Where("challenge_id = ?", challenge.ID)
		flags, _, ret := db.InitTeamFlagRepo(tx.Where("challenge_flag_id IN (?)", ids)).List(-1, -1, db.GetOptions{Conditions: map[string]any{"team_id": teamID}, Sort: []string{"challenge_flag_id ASC"}})
		if !ret.OK {
			return "", ret
		}
		for _, flag := range flags {
			values = append(values, flag.Value)
		}
	}
	return challenge.AttachmentCachePath(teamID, values), model.SuccessRetVal()
}
