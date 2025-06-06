package service

import (
	"CBCTF/internel/model"
	db "CBCTF/internel/repo"
	"gorm.io/gorm"
)

// CheckIfGenerated contestChallenge 需要预加载 ContestFlags
func CheckIfGenerated(tx *gorm.DB, team model.Team, contestChallenge model.ContestChallenge) bool {
	teamFlagRepo := db.InitTeamFlagRepo(tx)
	for _, contestFlag := range contestChallenge.ContestFlags {
		if _, ok, _ := teamFlagRepo.GetWithConditions(db.GetOptions{
			{Key: "team_id", Value: team.ID, Op: "and"},
			{Key: "contest_flag_id", Value: contestFlag.ID, Op: "and"},
		}); !ok {
			return false
		}
	}
	return true
}
