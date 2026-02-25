package middleware

import (
	"CBCTF/internal/db"
	"CBCTF/internal/i18n"
	"CBCTF/internal/model"
	"CBCTF/internal/resp"
	"CBCTF/internal/service"

	"github.com/gin-gonic/gin"
)

func CheckTeamVictimCount(ctx *gin.Context) {
	contest := GetContest(ctx)
	team := GetTeam(ctx)
	count, ret := service.CountTeamVictims(db.DB, team)
	if !ret.OK {
		resp.AbortJSON(ctx, ret)
		return
	}
	if count >= contest.Victims {
		resp.AbortJSON(ctx, model.RetVal{Msg: i18n.Model.Victim.Limited})
		return
	}
	ctx.Next()
}
