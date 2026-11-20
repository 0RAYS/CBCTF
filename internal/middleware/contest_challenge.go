package middleware

import (
	"github.com/gin-gonic/gin"

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
	checkTeamFlagState(ctx, true, i18n.Model.TeamFlag.AlreadySolved, service.CheckIfSolved)
}
