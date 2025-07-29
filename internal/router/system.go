package router

import (
	"CBCTF/internal/config"
	"CBCTF/internal/i18n"
	"CBCTF/internal/middleware"
	"CBCTF/internal/redis"
	db "CBCTF/internal/repo"
	"CBCTF/internal/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/net"
	"net/http"
)

func HomePage(ctx *gin.Context) {
	DB := db.DB.WithContext(ctx)

	data := gin.H{
		"upcoming":   []gin.H{},
		"stats":      []gin.H{},
		"scoreboard": []gin.H{},
	}
	contests, count, ok, _ := db.InitContestRepo(DB).List(-1, -1, db.GetOptions{
		Preloads: map[string]db.GetOptions{"Teams": {}, "Users": {}},
	})
	if ok {
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
				"users":    len(contest.Users),
				"teams":    len(contest.Teams),
				"avatar":   contest.Avatar,
			}
			data["upcoming"] = append(data["upcoming"].([]gin.H), info)
		}
	}
	data["stats"] = append(data["stats"].([]gin.H), gin.H{"label": "CTF Events", "value": count})
	count, _, _ = db.InitUserRepo(DB).Count()
	data["stats"] = append(data["stats"].([]gin.H), gin.H{"label": "Activate CTFers", "value": count})
	count, _, _ = db.InitChallengeRepo(DB).Count()
	data["stats"] = append(data["stats"].([]gin.H), gin.H{"label": "Challenges", "value": count})
	count, _, _ = db.InitSubmissionRepo(DB).Count()
	data["stats"] = append(data["stats"].([]gin.H), gin.H{"label": "Submissions", "value": count})
	users, _, _, _ := service.GetUserRanking(DB, 5, 0)
	for _, user := range users {
		data["scoreboard"] = append(data["scoreboard"].([]gin.H), gin.H{
			"name":    user.Name,
			"score":   user.Score,
			"solved":  user.Solved,
			"country": user.Country,
		})
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": i18n.Success, "data": data})
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

	var DB = db.DB.WithContext(ctx)
	ret["users"], _, _ = db.InitUserRepo(DB).Count()
	ret["contests"], _, _ = db.InitContestRepo(DB).Count()
	ret["ip"], _, _ = db.InitRequestRepo(DB).CountIP()
	ret["challenges"], _, _ = db.InitChallengeRepo(DB).Count()
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
	ctx.JSON(http.StatusOK, gin.H{"msg": i18n.Success, "data": ret})
}

func SystemConfig(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"msg": i18n.Success, "data": config.Env})
}
