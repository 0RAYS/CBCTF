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

func CheckTeamVictimCount(ctx *gin.Context) {
	contest := GetContest(ctx)
	team := GetTeam(ctx)
	// 靶场实例按池领取,不受每队靶机数量限制。
	if challenge := GetChallenge(ctx); challenge.RangeConfig != nil {
		ctx.Next()
		return
	}
	count, ret := service.CountTeamVictims(db.DB.WithContext(ctx.Request.Context()), team)
	if !ret.OK {
		resp.AbortJSON(ctx, ret)
		return
	}
	if count >= contest.Victims {
		resp.AbortJSON(ctx, model.FailRetVal(i18n.Model.Victim.Limited, "Victim", fmt.Sprintf("team_id=%d,limit=%d", team.ID, contest.Victims), nil))
		return
	}
	ctx.Next()
}
