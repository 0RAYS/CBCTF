package resp

import (
	"CBCTF/internel/model"
	"github.com/gin-gonic/gin"
)

func GetVictimResp(victim model.Victim) gin.H {
	return gin.H{
		"id":                   victim.ID,
		"team_id":              victim.TeamID,
		"contest_challenge_id": victim.ContestChallengeID,
		"ip_block":             victim.IPBlock,
		"start":                victim.Start,
		"duration":             victim.Duration,
	}
}
