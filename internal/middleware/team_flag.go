package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"CBCTF/internal/db"
	"CBCTF/internal/i18n"
	"CBCTF/internal/model"
	"CBCTF/internal/resp"
	"CBCTF/internal/service"
)

// checkTeamFlagState aborts the request when the team flag state of the
// current challenge differs from the expected one.
func checkTeamFlagState(
	ctx *gin.Context,
	expected bool,
	msg string,
	check func(*gorm.DB, model.Team, []model.ChallengeFlag) (bool, model.RetVal),
) {
	team := GetTeam(ctx)
	challenge := GetChallenge(ctx)
	flags, _, ret := db.InitChallengeFlagRepo(db.DB.WithContext(ctx.Request.Context())).List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"challenge_id": challenge.ID},
	})
	if !ret.OK {
		resp.AbortJSON(ctx, ret)
		return
	}
	state, ret := check(db.DB.WithContext(ctx.Request.Context()), team, flags)
	if !ret.OK {
		resp.AbortJSON(ctx, ret)
		return
	}
	if state == expected {
		resp.AbortJSON(ctx, model.FailRetVal(msg, "TeamFlag", fmt.Sprintf("team_id=%d,challenge_id=%d", team.ID, challenge.ID), nil))
		return
	}
	ctx.Next()
}

// CheckIfGenerated model.Team 是否初始化 model.TeamFlag
func CheckIfGenerated(ctx *gin.Context) {
	checkTeamFlagState(ctx, false, i18n.Model.TeamFlag.NotFound, service.CheckIfGenerated)
}
