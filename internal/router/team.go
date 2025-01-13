package router

import (
	"RayWar/internal/db"
	"RayWar/internal/log"
	"RayWar/internal/middleware"
	"RayWar/internal/model"
	"RayWar/internal/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetTeams(ctx *gin.Context) {
	var getTeamsForm GetTeamsForm
	var (
		teams []model.Team
		total int64
		ok    bool
		msg   string
	)
	trace := middleware.GetTraceID(ctx)
	if err := ctx.ShouldBind(&getTeamsForm); err == nil {
		if middleware.GetSelf(ctx).Type == "admin" {
			teams, total, ok, msg = db.GetTeams(getTeamsForm.Limit, getTeamsForm.Offset, true)
		} else {
			teams, total, ok, msg = db.GetTeams(getTeamsForm.Limit, getTeamsForm.Offset, false, getTeamsForm.ContestID)
		}
		if !ok {
			log.Logger.Infof("| %s | %s", trace, msg)
			ctx.JSONP(http.StatusInternalServerError, gin.H{"trace": trace, "msg": utils.M(ctx, msg), "data": nil})
			return
		}
		ctx.JSONP(http.StatusOK, gin.H{"trace": trace, "msg": utils.M(ctx, msg), "data": gin.H{
			"teams": utils.TidyRetData(teams, "captcha"),
			"total": total,
		}})
		return
	} else {
		log.Logger.Infof("| %s | %s", trace, err)
		ctx.JSONP(http.StatusBadRequest, gin.H{"trace": trace, "msg": utils.M(ctx, "BadRequest"), "data": nil})
		return
	}
}
