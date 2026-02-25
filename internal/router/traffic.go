package router

import (
	"CBCTF/internal/dto"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
	"CBCTF/internal/resp"
	"CBCTF/internal/service"
	"sort"

	"github.com/gin-gonic/gin"
)

func GetTraffics(ctx *gin.Context) {
	var form dto.GetTrafficForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	victim := middleware.GetVictim(ctx)
	connections, ipL, totalDuration, ret := service.GetTraffic(victim, form)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	sort.Strings(ipL)
	data := resp.GetTrafficResp(connections)
	data["ip"] = ipL
	data["duration"] = totalDuration
	resp.JSON(ctx, model.SuccessRetVal(data))
}
