package middleware

import (
	"RayWar/internal/model"
	"github.com/gin-gonic/gin"
)

// GetTraceID 从 gin.Context 中获取 trace，该值由 middleware.Trace 设置
func GetTraceID(ctx *gin.Context) string {
	return ctx.GetString("trace")
}

// GetSelf 从 gin.Context 中获取 self(model.User)，该值由 middleware.CheckLogin 设置
func GetSelf(ctx *gin.Context) model.User {
	self, ok := ctx.Get("self")
	if !ok {
		return model.User{}
	}
	return self.(model.User)
}

// GetUser 从 gin.Context 中获取 user(model.User)，该值由 middleware.SetUser 设置
func GetUser(ctx *gin.Context) model.User {
	user, ok := ctx.Get("user")
	if !ok {
		return model.User{}
	}
	return user.(model.User)
}

// GetTeam 从 gin.Context 中获取 team(model.User)，该值由 middleware.SetTeam 设置
func GetTeam(ctx *gin.Context) model.Team {
	team, ok := ctx.Get("team")
	if !ok {
		return model.Team{}
	}
	return team.(model.Team)
}

// GetContest 从 gin.Context 中获取 contest(model.Contest)，该值由 middleware.SetContest 设置
func GetContest(ctx *gin.Context) model.Contest {
	contest, ok := ctx.Get("contest")
	if !ok {
		return model.Contest{}
	}
	return contest.(model.Contest)
}
