package middleware

import (
	"CBCTF/internal/db"
	"CBCTF/internal/i18n"
	"CBCTF/internal/model"
	"CBCTF/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CheckTeamVictimCount(ctx *gin.Context) {
	contest := GetContest(ctx)
	team := GetTeam(ctx)
	count, ret := service.CountTeamVictims(db.DB, team)
	if !ret.OK {
		ctx.AbortWithStatusJSON(http.StatusOK, ret)
		return
	}
	if count >= contest.Victims {
		ctx.AbortWithStatusJSON(http.StatusOK, model.RetVal{Msg: i18n.Model.Victim.Limited})
		return
	}
	ctx.Next()
}
