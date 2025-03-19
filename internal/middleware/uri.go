package middleware

import (
	"CBCTF/internal/db"
	"CBCTF/internal/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

// SetUser 保存 model.User 至上下文
func SetUser(ctx *gin.Context) {
	type userIDUri struct {
		UserID uint `uri:"userID" binding:"required"`
	}
	var userID userIDUri
	if err := ctx.ShouldBindUri(&userID); err != nil {
		ctx.JSONP(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		ctx.Abort()
		return
	}
	user, ok, msg := db.GetUserByID(db.DB.WithContext(ctx), userID.UserID)
	if !ok {
		ctx.JSONP(http.StatusOK, gin.H{"msg": msg, "data": nil})
		ctx.Abort()
		return
	}
	ctx.Set("User", user)
	ctx.Next()
}

// GetUser 从上下文中获取 model.User
func GetUser(ctx *gin.Context) model.User {
	if user, ok := ctx.Get("User"); !ok {
		return model.User{}
	} else {
		return user.(model.User)
	}
}

// SetContest 保存 model.Contest 至上下文
func SetContest(ctx *gin.Context) {
	type contestIDUri struct {
		ContestID uint `uri:"contestID" binding:"required"`
	}
	var contestID contestIDUri
	if err := ctx.ShouldBindUri(&contestID); err != nil {
		ctx.JSONP(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		ctx.Abort()
		return
	}
	contest, ok, msg := db.GetContestByID(db.DB.WithContext(ctx), contestID.ContestID)
	if !ok {
		ctx.JSONP(http.StatusOK, gin.H{"msg": msg, "data": nil})
		ctx.Abort()
		return
	}
	ctx.Set("Contest", contest)
	ctx.Next()
}

// GetContest 从上下文中获取 model.Contest
func GetContest(ctx *gin.Context) model.Contest {
	if contest, ok := ctx.Get("Contest"); !ok {
		return model.Contest{}
	} else {
		return contest.(model.Contest)
	}
}

// SetTeamByUser 依照 model.User model.Contest 保存 model.Team 至上下文, 调用前前文须设置 model.Contest
func SetTeamByUser(ctx *gin.Context) {
	var (
		self model.User
		team model.Team
		ok   bool
		msg  string
	)
	self, ok = GetSelf(ctx).(model.User)
	if !ok {
		ctx.JSONP(http.StatusForbidden, gin.H{"msg": "Forbidden", "data": nil})
		ctx.Abort()
		return
	}
	team, ok, msg = db.GetTeamByUserID(db.DB.WithContext(ctx), self.ID, GetContest(ctx).ID)
	if !ok {
		ctx.JSONP(http.StatusOK, gin.H{"msg": msg, "data": nil})
		ctx.Abort()
		return
	}
	ctx.Set("Team", team)
	ctx.Next()
}

// SetTeamByURI 保存 model.Team 至上下文
func SetTeamByURI(ctx *gin.Context) {
	var (
		team model.Team
		ok   bool
		msg  string
	)
	type teamIDUri struct {
		TeamID uint `uri:"teamID" binding:"required"`
	}
	var teamID teamIDUri
	if err := ctx.ShouldBindUri(&teamID); err != nil {
		ctx.JSONP(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		ctx.Abort()
		return
	}
	team, ok, msg = db.GetTeamByID(db.DB.WithContext(ctx), teamID.TeamID)
	if !ok {
		ctx.JSONP(http.StatusOK, gin.H{"msg": msg, "data": nil})
		ctx.Abort()
		return
	}
	ctx.Set("Team", team)
	ctx.Next()
}

// GetTeam 从上下文中获取 model.Team
func GetTeam(ctx *gin.Context) model.Team {
	if team, ok := ctx.Get("Team"); !ok {
		return model.Team{}
	} else {
		return team.(model.Team)
	}
}

// SetFile 保存 model.File 至上下文
func SetFile(t string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		type fileIDUri struct {
			FileID string `uri:"fileID" binding:"required"`
		}
		var fileID fileIDUri
		if err := ctx.ShouldBindUri(&fileID); err != nil {
			ctx.JSONP(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
			ctx.Abort()
			return
		}
		file, ok, msg := db.GetFileByID(db.DB.WithContext(ctx), fileID.FileID)
		if !ok {
			ctx.JSONP(http.StatusOK, gin.H{"msg": msg, "data": nil})
			ctx.Abort()
			return
		}
		if file.Type != t {
			ctx.JSONP(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
			ctx.Abort()
			return
		}
		ctx.Set("File", file)
		ctx.Next()
	}
}

// GetFile 从上下文中获取 model.File
func GetFile(ctx *gin.Context) model.File {
	if file, ok := ctx.Get("File"); !ok {
		return model.File{}
	} else {
		return file.(model.File)
	}
}

// SetChallenge 保存 model.Challenge 至上下文
func SetChallenge(ctx *gin.Context) {
	type challengeIDUri struct {
		ChallengeID string `uri:"challengeID" binding:"required"`
	}
	var challengeID challengeIDUri
	if err := ctx.ShouldBindUri(&challengeID); err != nil {
		ctx.JSONP(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		ctx.Abort()
		return
	}
	challenge, ok, msg := db.GetChallengeByID(db.DB.WithContext(ctx), challengeID.ChallengeID)
	if !ok {
		ctx.JSONP(http.StatusOK, gin.H{"msg": msg, "data": nil})
		ctx.Abort()
		return
	}
	ctx.Set("Challenge", challenge)
	ctx.Next()
}

// GetChallenge 从上下文中获取 model.Challenge
func GetChallenge(ctx *gin.Context) model.Challenge {
	if challenge, ok := ctx.Get("Challenge"); !ok {
		return model.Challenge{}
	} else {
		return challenge.(model.Challenge)
	}
}

// SetUsage 保存 model.Usage 至上下文
func SetUsage(ctx *gin.Context) {
	type challengeIDUri struct {
		ChallengeID string `uri:"challengeID" binding:"required"`
	}
	var challengeID challengeIDUri
	if err := ctx.ShouldBindUri(&challengeID); err != nil {
		ctx.JSONP(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		ctx.Abort()
		return
	}
	usage, ok, msg := db.GetUsageBy2ID(db.DB.WithContext(ctx), GetContest(ctx).ID, challengeID.ChallengeID)
	if !ok {
		ctx.JSONP(http.StatusOK, gin.H{"msg": msg, "data": nil})
		ctx.Abort()
		return
	}
	ctx.Set("Usage", usage)
	ctx.Next()
}

// GetUsage 从上下文中获取 model.Usage
func GetUsage(ctx *gin.Context) model.Usage {
	if usage, ok := ctx.Get("Usage"); !ok {
		return model.Usage{}
	} else {
		return usage.(model.Usage)
	}
}

// SetContainer 保存 model.Docker 至上下文
func SetContainer(ctx *gin.Context) {
	type containerIDUri struct {
		ContainerID uint `uri:"containerID" binding:"required"`
	}
	var containerID containerIDUri
	if err := ctx.ShouldBindUri(&containerID); err != nil {
		ctx.JSONP(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		ctx.Abort()
		return
	}
	container, ok, msg := db.GetDockerByID(db.DB.WithContext(ctx), containerID.ContainerID, true)
	if !ok {
		ctx.JSONP(http.StatusOK, gin.H{"msg": msg, "data": nil})
		ctx.Abort()
		return
	}
	ctx.Set("Container", container)
	ctx.Next()
}

// GetContainer 从上下文中获取 model.Docker
func GetContainer(ctx *gin.Context) model.Docker {
	if container, ok := ctx.Get("Container"); !ok {
		return model.Docker{}
	} else {
		return container.(model.Docker)
	}
}

// SetNotice 保存 model.Notice 至上下文
func SetNotice(ctx *gin.Context) {
	type noticeIDUri struct {
		NoticeID uint `uri:"noticeID" binding:"required"`
	}
	var noticeID noticeIDUri
	if err := ctx.ShouldBindUri(&noticeID); err != nil {
		ctx.JSONP(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		ctx.Abort()
		return
	}
	notice, ok, msg := db.GetNoticeByID(db.DB.WithContext(ctx), noticeID.NoticeID)
	if !ok {
		ctx.JSONP(http.StatusOK, gin.H{"msg": msg, "data": nil})
		ctx.Abort()
		return
	}
	ctx.Set("Notice", notice)
	ctx.Next()
}

// GetNotice 从上下文中获取 model.Notice
func GetNotice(ctx *gin.Context) model.Notice {
	if notice, ok := ctx.Get("Notice"); !ok {
		return model.Notice{}
	} else {
		return notice.(model.Notice)
	}
}
