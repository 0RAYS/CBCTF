package router

import (
	f "CBCTF/internal/form"
	"CBCTF/internal/i18n"
	"CBCTF/internal/middleware"
	r "CBCTF/internal/redis"
	"fmt"
	"net/http"
	"slices"
	"time"

	"github.com/gin-gonic/gin"
)

func GetTraffics(ctx *gin.Context) {
	var form f.GetTrafficForm
	if ok, msg := form.Bind(ctx); !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	victim := middleware.GetVictim(ctx)
	connections, ok, msg := r.GetTraffic(victim)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	if len(connections) < 1 {
		ok, msg = r.UpdateTraffics(victim)
		if !ok {
			ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
			return
		}
		connections, ok, msg = r.GetTraffic(victim)
		if !ok {
			ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
			return
		}
	}
	totalDuation := connections[len(connections)-1].Time.Sub(connections[0].Time).Seconds()
	firstPacket := connections[0]
	firstPacket.TimeShift = 0
	startIndex := 0
	endIndex := len(connections) - 1
	for i, connection := range connections {
		connections[i].TimeShift = connection.Time.Sub(firstPacket.Time)
		if connections[i].TimeShift > time.Duration(form.TimeShift*1e9) {
			startIndex = i - 1
			if startIndex < 0 {
				startIndex = 0
			}
		}
		if connections[i].TimeShift > time.Duration((form.TimeShift+form.Duration)*1e9) {
			endIndex = i
			break
		}
	}
	type Statistics struct {
		SrcIP   string
		DstIP   string
		SrcPort uint16
		DstPort uint16
		Type    string
		Count   int64
		Size    int64
	}
	data := make(map[string]Statistics)
	ipL := make([]string, 0)
	for _, connection := range connections[startIndex:endIndex] {
		key := fmt.Sprintf("%s:%d-%s:%d-%s", connection.SrcIP, connection.SrcPort, connection.DstIP, connection.DstPort, connection.Type)
		if stats, exists := data[key]; exists {
			stats.Count += 1
			stats.Size += int64(connection.Size)
			data[key] = stats
		} else {
			data[key] = Statistics{
				SrcIP:   connection.SrcIP,
				DstIP:   connection.DstIP,
				SrcPort: connection.SrcPort,
				DstPort: connection.DstPort,
				Count:   1,
				Size:    int64(connection.Size),
			}
		}
		if !slices.Contains(ipL, connection.SrcIP) {
			ipL = append(ipL, connection.SrcIP)
		}
		if !slices.Contains(ipL, connection.DstIP) {
			ipL = append(ipL, connection.DstIP)
		}
	}
	conn := make([]gin.H, 0)
	for _, stats := range data {
		conn = append(conn, gin.H{
			"src_ip":   stats.SrcIP,
			"dst_ip":   stats.DstIP,
			"src_port": stats.SrcPort,
			"dst_port": stats.DstPort,
			"type":     stats.Type,
			"count":    stats.Count,
			"size":     stats.Size,
		})
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": i18n.Success, "data": gin.H{"connections": conn, "ip": ipL, "duration": int64(totalDuation)/1e9 + 1}})
}
