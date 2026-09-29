package middleware

import (
	"github.com/gin-gonic/gin"

	"CBCTF/internal/db"
	"CBCTF/internal/i18n"
	"CBCTF/internal/model"
	"CBCTF/internal/resp"
	"CBCTF/internal/service"
)

func CheckChallengeType(t model.ChallengeType) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if GetChallenge(ctx).Type != t {
			resp.AbortJSON(ctx, model.RetVal{Msg: i18n.Model.Challenge.InvalidType})
			return
		}
		ctx.Next()
	}
}

// CheckSolved model.Team 是否完全解决 model.ContestChallenge
func CheckSolved(ctx *gin.Context) {
	team := GetTeam(ctx)
	contestChallenge := GetContestChallenge(ctx)
	contestFlags, _, ret := db.InitContestFlagRepo(db.DB.WithContext(ctx.Request.Context())).List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"contest_challenge_id": contestChallenge.ID},
	})
	if !ret.OK {
		resp.AbortJSON(ctx, ret)
		return
	}
	solved, ret := service.CheckIfSolved(db.DB.WithContext(ctx.Request.Context()), team, contestFlags)
	if !ret.OK {
		resp.AbortJSON(ctx, ret)
		return
	}
	if solved {
		resp.AbortJSON(ctx, model.RetVal{Msg: i18n.Model.TeamFlag.AlreadySolved})
		return
	}
	ctx.Next()
}
