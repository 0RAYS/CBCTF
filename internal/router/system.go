package router

import (
	"CBCTF/internal/config"
	"CBCTF/internal/db"
	"CBCTF/internal/redis"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"net/http"
)

func SystemStatus(ctx *gin.Context) {
	ret := make(map[string]interface{})
	cpuPercent, err := cpu.Percent(0, false)
	if err != nil || len(cpuPercent) == 0 {
		ret["cpu"] = "0.00"
	} else {
		ret["cpu"] = fmt.Sprintf("%.2f", cpuPercent[0])
	}

	vm, err := mem.VirtualMemory()
	if err != nil || vm == nil {
		ret["mem"] = "0.00"
	} else {
		ret["mem"] = fmt.Sprintf("%.2f", vm.UsedPercent)
	}

	diskStat, err := disk.Usage("/") // 根目录
	if err != nil || diskStat == nil {
		ret["disk"] = "0.00"
	} else {
		ret["disk"] = fmt.Sprintf("%.2f", diskStat.UsedPercent)
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
