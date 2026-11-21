package resp

import (
	"github.com/gin-gonic/gin"

	"CBCTF/internal/model"
	"CBCTF/internal/view"
)

func GetVictimStatusResp(status view.VictimStatusView) gin.H {
	return gin.H{
		"target":    status.Targets,
		"duration":  status.Duration,
		"remaining": status.Remaining,
		"status":    status.Status,
	}
}

func GetVictimResp(victim model.Victim) gin.H {
	return gin.H{
		"id":                     victim.ID,
		"team_id":                victim.TeamID.V,
		"contest_id":             victim.ContestID.V,
		"challenge_id":           victim.ChallengeID,
		"challenge_name":         victim.Challenge.Name,
		"user_id":                victim.UserID,
		"start":                  victim.Start,
		"duration":               int64(victim.Duration.Seconds()),
		"status":                 victim.Status,
	}
}
