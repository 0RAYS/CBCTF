package router

import (
	f "CBCTF/internal/form"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
	"CBCTF/internal/resp"
	"CBCTF/internal/service"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
)

func GetTraffics(ctx *gin.Context) {
	var form f.GetTrafficForm
	if ret := form.Bind(ctx); !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	victim := middleware.GetVictim(ctx)
	connections, ipL, totalDuration, ret := service.GetTraffic(victim, form)
	if !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	sort.Strings(ipL)
	data := resp.GetTrafficResp(connections)
	data["ip"] = ipL
	data["duration"] = totalDuration
	ctx.JSON(http.StatusOK, model.SuccessRetVal(data))
}
