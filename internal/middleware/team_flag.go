package middleware

import (
	"github.com/gin-gonic/gin"

	"CBCTF/internal/db"
	"CBCTF/internal/i18n"
	"CBCTF/internal/model"
	"CBCTF/internal/resp"
	"CBCTF/internal/service"
)

// CheckIfGenerated model.Team 是否初始化 model.TeamFlag
func CheckIfGenerated(ctx *gin.Context) {
	team := GetTeam(ctx)
	contestChallenge := GetContestChallenge(ctx)
	contestFlags, _, ret := db.InitContestFlagRepo(db.DB.WithContext(ctx.Request.Context())).List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"contest_challenge_id": contestChallenge.ID},
	})
	if !ret.OK {
		resp.AbortJSON(ctx, ret)
		return
	}
	generated, ret := service.CheckIfGenerated(db.DB.WithContext(ctx.Request.Context()), team, contestFlags)
	if !ret.OK {
		resp.AbortJSON(ctx, ret)
		return
	}
	if !generated {
		resp.AbortJSON(ctx, model.RetVal{Msg: i18n.Model.TeamFlag.NotFound})
		return
	}
	ctx.Next()
}
