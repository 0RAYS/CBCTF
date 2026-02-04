package router

import (
	"CBCTF/internal/config"
	"CBCTF/internal/db"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
	"CBCTF/internal/prometheus"
	"CBCTF/internal/redis"
	"CBCTF/internal/service"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/net"
)

func HomePage(ctx *gin.Context) {
	data := gin.H{
		"upcoming":   []gin.H{},
		"stats":      []gin.H{},
		"scoreboard": []gin.H{},
	}
	repo := db.InitContestRepo(db.DB)
	contests, count, ret := repo.List(-1, -1)
	if ret.OK {
		for i := 0; i < func() int {
			if len(contests) > 3 {
				return 3
			}
			return len(contests)
		}(); i++ {
			contest := contests[i]
			info := gin.H{
				"name":     contest.Name,
				"start":    contest.Start,
				"duration": contest.Duration.Seconds(),
				"users":    repo.CountAssociation(contest, "Users"),
				"teams":    repo.CountAssociation(contest, "Teams"),
				"picture":  contest.Picture,
			}
			data["upcoming"] = append(data["upcoming"].([]gin.H), info)
		}
	}
	data["stats"] = append(data["stats"].([]gin.H), gin.H{"label": "CTF Events", "value": count})
	count, _ = db.InitUserRepo(db.DB).Count()
	data["stats"] = append(data["stats"].([]gin.H), gin.H{"label": "Activate CTFers", "value": count})
	count, _ = db.InitChallengeRepo(db.DB).Count()
	data["stats"] = append(data["stats"].([]gin.H), gin.H{"label": "Challenges", "value": count})
	count, _ = db.InitSubmissionRepo(db.DB).Count()
	data["stats"] = append(data["stats"].([]gin.H), gin.H{"label": "Submissions", "value": count})
	users, _, _ := service.GetUserRanking(db.DB, 5, 0)
	for _, user := range users {
		data["scoreboard"] = append(data["scoreboard"].([]gin.H), gin.H{
			"name":   user.Name,
			"score":  user.Score,
			"solved": user.Solved,
		})
	}
	ctx.JSON(http.StatusOK, model.SuccessRetVal(data))
}

func SystemStatus(ctx *gin.Context) {
	ret := make(map[string]any)
	ret["metrics"] = redis.GetMetrics()

	ioStats, err := net.IOCounters(false)
	if err != nil || len(ioStats) == 0 {
		ret["io"] = 0
		ret["sent"] = 0
		ret["recv"] = 0
	} else {
		ret["io"] = ioStats[0].BytesSent + ioStats[0].BytesRecv
		ret["sent"] = ioStats[0].BytesSent
		ret["recv"] = ioStats[0].BytesRecv
	}

	ret["users"], _ = db.InitUserRepo(db.DB).Count()
	ret["contests"], _ = db.InitContestRepo(db.DB).Count()
	ret["ip"], _ = db.InitRequestRepo(db.DB).CountIP()
	ret["challenges"], _ = db.InitChallengeRepo(db.DB).Count()
	middleware.MU.Lock()
	if middleware.TotalRequests == 0 {
		ret["requests"] = 0
		ret["duration"] = 0
	} else {
		ret["requests"] = middleware.TotalRequests
		ret["duration"] = middleware.TotalDuration.Milliseconds() / int64(middleware.TotalRequests)
	}
	middleware.MU.Unlock()

	total, hit, miss := redis.Status()
	ret["cache"] = total
	ret["hit"] = hit
	if hit+miss == 0 {
		ret["rate"] = "0.00"
	} else {
		ret["rate"] = fmt.Sprintf("%.2f", float64(hit)/float64(hit+miss)*100)
	}
	prometheus.UpdateCacheMetrics("redis", hit, miss)
	ctx.JSON(http.StatusOK, model.SuccessRetVal(ret))
}

func SystemConfig(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, model.SuccessRetVal(config.Env))
}
