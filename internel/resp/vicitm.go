package resp

import (
	"CBCTF/internel/model"
	"github.com/gin-gonic/gin"
)

func GetVictimResp(victim model.Victim) gin.H {
	return gin.H{
		"id":       victim.ID,
		"team_id":  victim.TeamID,
		"usage_id": victim.UsageID,
		"ip_block": victim.IPBlock,
		"start":    victim.Start,
		"duration": victim.Duration,
	}
}
