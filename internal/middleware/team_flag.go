package middleware

import (
	"CBCTF/internal/db"
	"CBCTF/internal/i18n"
	"CBCTF/internal/model"
	"CBCTF/internal/service"
	"net/http"

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
		ctx.JSON(http.StatusOK, ret)
		return
	}
	if !service.CheckIfGenerated(db.DB, team, contestFlags) {
		ctx.AbortWithStatusJSON(http.StatusOK, model.RetVal{Msg: i18n.Model.NotFound, Attr: map[string]any{"Model": model.TeamFlag{}.ModelName()}})
		return
	}
	ctx.Next()
}
