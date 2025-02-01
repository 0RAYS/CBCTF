package router

import (
	"CBCTF/internal/config"
	"CBCTF/internal/db"
	"CBCTF/internal/log"
	"CBCTF/internal/middleware"
	"CBCTF/internal/redis"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/net"
	"net/http"
	"reflect"
	"time"
)

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

	ret["users"] = db.CountUsers(ctx)
	ret["contests"] = db.CountContests(ctx)
	ret["ip"] = db.CountIP(ctx)
	ret["challenges"] = db.CountChallenges(ctx)
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

func SystemUpdate(ctx *gin.Context) {
	var env config.Config
	if err := ctx.ShouldBind(&env); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	if reflect.DeepEqual(env, *config.Env) {
		log.Logger.Debug("Config not change")
		ctx.JSON(http.StatusOK, gin.H{"msg": "ConfigNotChange", "data": nil})
		return
	}
	go func() {
		time.Sleep(time.Second * 2)
		err := config.Save(env)
		if err != nil {
			log.Logger.Warningf("Failed to save config: %s", err)
		}
	}()
	ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": nil})
}
