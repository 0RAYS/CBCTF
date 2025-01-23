package router

import (
	"CBCTF/internal/config"
	"CBCTF/internal/db"
	"CBCTF/internal/redis"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/net"
	"net/http"
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
