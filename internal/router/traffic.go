package router

import (
	f "CBCTF/internal/form"
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/middleware"
	"CBCTF/internal/redis"
	"CBCTF/internal/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetTraffics(ctx *gin.Context) {
	var form f.GetTrafficForm
	if ok, msg := form.Bind(ctx); !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	victim := middleware.GetVictim(ctx)
	connections, err := redis.GetTraffics(victim)
	if err != nil {
		log.Logger.Warningf("Failed to get traffics: %s", err)
		ctx.JSON(http.StatusOK, gin.H{"msg": i18n.RedisError, "data": nil})
		return
	}
	data := make([]utils.Connection, 0)
	for _, conn := range connections {
		if conn.Time.Unix() >= form.Start.Unix() && conn.Time.Unix() <= form.End.Unix() {
			data = append(data, conn)
		}
	}
	ctx.JSON(http.StatusOK, gin.H{"data": gin.H{"connections": data, "duration": form.End.Second() - form.Start.Second()}, "msg": i18n.Success})
}
