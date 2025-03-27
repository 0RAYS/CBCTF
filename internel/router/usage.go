package router

import (
	"CBCTF/internel/middleware"
	db "CBCTF/internel/repo"
	"CBCTF/internel/resp"
	"CBCTF/internel/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func GetUsages(ctx *gin.Context) {
	var (
		all     = middleware.GetRole(ctx) == "admin"
		DB      = db.DB.WithContext(ctx)
		contest = middleware.GetContest(ctx)
		team    = middleware.GetTeam(ctx)
	)
	usages, ok, msg := service.GetUsages(DB, contest, all)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	var data []gin.H
	for _, usage := range usages {
		tmp := resp.GetUsageResp(usage)
		tmp["attempts"], _, _ = service.CountAttempts(DB, team, usage)
		tmp["init"], _, _ = service.IsGenerated(DB, usage)
		tmp["solved"], _, _ = service.IsSolved(DB, team, usage)
		tmp["remote"] = service.GetRemoteStatus(DB, usage)
		tmp["file"] = func() string {
			if _, err := os.Stat(usage.Challenge.AttachmentPath(team.ID)); err != nil {
				return ""
			}
			return usage.Challenge.AttachmentPath(team.ID)
		}
		data = append(data, tmp)
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": data})
}
