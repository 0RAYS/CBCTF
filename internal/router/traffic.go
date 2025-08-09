package router

import (
	f "CBCTF/internal/form"
	"CBCTF/internal/i18n"
	"CBCTF/internal/middleware"
	"CBCTF/internal/resp"
	"CBCTF/internal/service"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
)

func GetTraffics(ctx *gin.Context) {
	var form f.GetTrafficForm
	if ok, msg := form.Bind(ctx); !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	victim := middleware.GetVictim(ctx)
	connections, ipL, totalDuration, ok, msg := service.GetTraffic(victim, form)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	sort.Strings(ipL)
	data := resp.GetTrafficResp(connections)
	data["ip"] = ipL
	data["duration"] = totalDuration
	ctx.JSON(http.StatusOK, gin.H{"msg": i18n.Success, "data": data})
}
