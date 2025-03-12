package router

import (
	"CBCTF/internal/config"
	"CBCTF/internal/db"
	"CBCTF/internal/middleware"
	"CBCTF/internal/redis"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/net"
	"net/http"
)

func HomePage(ctx *gin.Context) {
	data := gin.H{
		"upcoming":   []gin.H{},
		"stats":      []gin.H{},
		"scoreboard": []gin.H{},
	}
	DB := db.DB.WithContext(ctx)
	contests, contestCount, ok, _ := db.GetContests(DB, 3, 0, false, true, true)
	if ok {
		count := func() int {
			if len(contests) > 3 {
				return 3
			}
			return len(contests)
		}()
		if count > 0 {
			for i := 1; i < count; i++ {
				contest := contests[i]
				info := gin.H{
					"name":     contest.Name,
					"start":    contest.Start,
					"duration": contest.Duration.Seconds(),
					"users":    len(contest.Users),
					"avatar":   contest.Avatar,
				}
				data["upcoming"] = append(data["upcoming"].([]gin.H), info)
			}
		}
	}
	data["stats"] = append(data["stats"].([]gin.H), gin.H{"label": "CTF Events", "value": contestCount})
	_, userCount, _, _ := db.GetUsers(DB, 0, 0, true, false)
	data["stats"] = append(data["stats"].([]gin.H), gin.H{"label": "Activate CTFers", "value": userCount})
	_, challengeCount, _, _ := db.GetChallenges(DB, 0, 0, "", "")
	data["stats"] = append(data["stats"].([]gin.H), gin.H{"label": "Challenges", "value": challengeCount})
	_, submissionCount, _, _ := db.GetSubmissions(DB, 0, 0, "")
	data["stats"] = append(data["stats"].([]gin.H), gin.H{"label": "Submissions", "value": submissionCount})
	users, _, _, _ := db.GetUserRanking(DB, 5, 0)
	for _, user := range users {
		data["scoreboard"] = append(data["scoreboard"].([]gin.H), gin.H{
			"name":    user.Name,
			"score":   user.Score,
			"solved":  user.Solved,
			"country": user.Country,
		})
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": data})
}

func SystemStatus(ctx *gin.Context) {
	ret := make(map[string]interface{})
	metrics, err := redis.GetMetrics()
	if err != nil {
		ret["metrics"] = nil
	} else {
		ret["metrics"] = metrics
	}

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
	ret["users"] = db.CountUsers(DB)
	ret["contests"] = db.CountContests(DB)
	ret["ip"] = db.CountIP(DB)
	ret["challenges"] = db.CountChallenges(DB)
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
	ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": ret})
}

func SystemConfig(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": config.Env})
}

//func SystemUpdate(ctx *gin.Context) {
//	var env config.Config
//	if err := ctx.ShouldBind(&env); err != nil {
//		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
//		return
//	}
//	if reflect.DeepEqual(env, *config.Env) {
//		log.Logger.Debug("Config not change")
//		ctx.JSON(http.StatusOK, gin.H{"msg": "ConfigNotChange", "data": nil})
//		return
//	}
//	go func() {
//		time.Sleep(time.Second * 2)
//		err := config.Save(env)
//		if err != nil {
//			log.Logger.Warningf("Failed to save config: %s", err)
//		}
//	}()
//	ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": nil})
//}
