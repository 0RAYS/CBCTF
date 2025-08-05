package router

import (
	f "CBCTF/internal/form"
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/middleware"
	"CBCTF/internal/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strings"
)

func GetTraffics(ctx *gin.Context) {
	var form f.GetTrafficForm
	if ok, msg := form.Bind(ctx); !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	victim := middleware.GetVictim(ctx)
	dir, err := os.ReadDir(victim.TrafficBasePath())
	if err != nil {
		log.Logger.Warningf("Failed to read dir: %v", err)
		ctx.JSON(http.StatusOK, gin.H{"msg": i18n.UnknownError, "data": nil})
		return
	}
	connections := make([]utils.Connection, 0)
	for _, file := range dir {
		if file.IsDir() || (!strings.HasSuffix(file.Name(), ".pcap") && !strings.HasSuffix(file.Name(), ".pcapng")) {
			continue
		}
		packet, err := utils.ReadPcap(fmt.Sprintf("%s/%s", victim.TrafficBasePath(), file.Name()))
		if err != nil {
			ctx.JSON(http.StatusOK, gin.H{"msg": i18n.UnknownError, "data": nil})
			return
		}
		connections = append(connections, packet...)
	}
	data := make([]utils.Connection, 0)
	for _, conn := range connections {
		if conn.Time.Unix() >= form.Start.Unix() && conn.Time.Unix() <= form.End.Unix() {
			data = append(data, conn)
		}
	}
	ctx.JSON(http.StatusOK, gin.H{"data": gin.H{"connections": data, "duration": form.End.Second() - form.Start.Second()}, "msg": i18n.Success})
}
