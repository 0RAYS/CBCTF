package middleware

import (
	"CBCTF/internal/db"
	"CBCTF/internal/i18n"
	"CBCTF/internal/model"
	"CBCTF/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CheckChallengeType(t string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if GetChallenge(ctx).Type != t {
			ctx.AbortWithStatusJSON(http.StatusOK, model.RetVal{Msg: i18n.Model.Challenge.InvalidType})
			return
		}
		ctx.Next()
	}
}

// CheckSolved model.Team 是否完全解决 model.ContestChallenge
func CheckSolved(ctx *gin.Context) {
	team := GetTeam(ctx)
	contestChallenge := GetContestChallenge(ctx)
	contestFlags, _, ret := db.InitContestFlagRepo(db.DB).List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"contest_challenge_id": contestChallenge.ID},
	})
	if !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	if service.CheckIfSolved(db.DB, team, contestFlags) {
		ctx.AbortWithStatusJSON(http.StatusOK, model.RetVal{Msg: i18n.Model.TeamFlag.AlreadySolved})
		return
	}
	ctx.Next()
}
