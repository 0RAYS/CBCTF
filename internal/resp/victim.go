package resp

import (
	"CBCTF/internal/model"

	"github.com/gin-gonic/gin"
)

func GetVictimResp(victim model.Victim) gin.H {
	return gin.H{
		"id":                     victim.ID,
		"team_id":                victim.TeamID.V,
		"challenge_id":           victim.ChallengeID,
		"contest_id":             victim.ContestID.V,
		"contest_challenge_id":   victim.ContestChallengeID.V,
		"contest_challenge_name": victim.ContestChallenge.Name,
		"user_id":                victim.UserID.V,
		"start":                  victim.Start,
		"duration":               victim.Duration,
	}
}
