package resp

import (
	"github.com/gin-gonic/gin"

	"CBCTF/internal/view"
)

func GetRangePoolResp(pool view.RangePoolView) gin.H {
	instances := make([]gin.H, 0, len(pool.Instances))
	for _, instance := range pool.Instances {
		instances = append(instances, gin.H{
			"id":              instance.ID,
			"status":          instance.Status,
			"team_id":         instance.TeamID,
			"team_name":       instance.TeamName,
			"lease_expire_at": instance.LeaseExpireAt,
			"template_key":    instance.TemplateKey,
			"nodes":           instance.Nodes,
			"start":           instance.Start,
		})
	}
	return gin.H{
		"challenge_id":   pool.ChallengeID,
		"challenge_name": pool.ChallengeName,
		"pool_size":      pool.PoolSize,
		"max_leased":     pool.MaxLeased,
		"leased":         pool.Leased,
		"instances":      instances,
	}
}

func GetRangePoolsResp(pools []view.RangePoolView) gin.H {
	data := make([]gin.H, 0, len(pools))
	for _, pool := range pools {
		data = append(data, GetRangePoolResp(pool))
	}
	return gin.H{"pools": data}
}
