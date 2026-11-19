package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"CBCTF/internal/db"
	"CBCTF/internal/i18n"
	"CBCTF/internal/model"
	"CBCTF/internal/resp"
	"CBCTF/internal/service"
)

func CheckChallengeType(t model.ChallengeType) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		challengeType := GetChallenge(ctx).Type
		if contestChallenge := GetContestChallenge(ctx); contestChallenge.ID > 0 {
			challengeType = contestChallenge.Type
		}
		if challengeType != t {
			resp.AbortJSON(ctx, model.FailRetVal(i18n.Model.Challenge.InvalidType, "Challenge", "type="+string(challengeType), nil))
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
		resp.AbortJSON(ctx, model.FailRetVal(i18n.Model.TeamFlag.AlreadySolved, "TeamFlag", fmt.Sprintf("team_id=%d,contest_challenge_id=%d", team.ID, contestChallenge.ID), nil))
		return
	}
	ctx.Next()
}
