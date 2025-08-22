package resp

import (
	"CBCTF/internal/model"

	"github.com/gin-gonic/gin"
)

func GetVictimResp(victim model.Victim) gin.H {
	return gin.H{
		"id":                   victim.ID,
		"team_id":              victim.TeamID,
		"challenge_id":         victim.ChallengeID,
		"contest_challenge_id": victim.ContestChallengeID,
		"start":                victim.Start,
		"duration":             victim.Duration,
	}
}
