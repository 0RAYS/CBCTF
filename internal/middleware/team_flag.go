package middleware

import (
	"CBCTF/internal/db"
	"CBCTF/internal/i18n"
	"CBCTF/internal/model"
	"CBCTF/internal/resp"
	"CBCTF/internal/service"

	"github.com/gin-gonic/gin"
)

// CheckIfGenerated model.Team 是否初始化 model.TeamFlag
func CheckIfGenerated(ctx *gin.Context) {
	team := GetTeam(ctx)
	contestChallenge := GetContestChallenge(ctx)
	contestFlags, _, ret := db.InitContestFlagRepo(db.DB).List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"contest_challenge_id": contestChallenge.ID},
	})
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	if !service.CheckIfGenerated(db.DB, team, contestFlags) {
		resp.AbortJSON(ctx, model.RetVal{Msg: i18n.Model.TeamFlag.NotFound})
		return
	}
	ctx.Next()
}
