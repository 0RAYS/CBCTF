package middleware

import (
	"CBCTF/internel/model"
	db "CBCTF/internel/repo"
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
	user, ok, msg := db.InitUserRepo(db.DB.WithContext(ctx)).GetByID(userID.UserID, "all")
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
	contest, ok, msg := db.InitContestRepo(db.DB.WithContext(ctx)).GetByID(contestID.ContestID, "all")
	if !ok {
		ctx.JSONP(http.StatusOK, gin.H{"msg": msg, "data": nil})
		ctx.Abort()
		return
	}
	if GetRole(ctx) != "admin" && contest.Hidden {
		ctx.JSONP(http.StatusNotFound, gin.H{"msg": "ContestNotFound", "data": nil})
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
	team, ok, msg = db.InitTeamRepo(db.DB.WithContext(ctx)).GetBy2ID(self.ID, GetContest(ctx).ID)
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
	team, ok, msg = db.InitTeamRepo(db.DB.WithContext(ctx)).GetByID(teamID.TeamID, "all")
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
		file, ok, msg := db.InitFileRepo(db.DB.WithContext(ctx)).GetByID(fileID.FileID)
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
	challenge, ok, msg := db.InitChallengeRepo(db.DB.WithContext(ctx)).GetByID(challengeID.ChallengeID, "all")
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
	usage, ok, msg := db.InitUsageRepo(db.DB.WithContext(ctx)).GetBy2ID(GetContest(ctx).ID, challengeID.ChallengeID, true, "all")
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

func SetFlag(ctx *gin.Context) {
	type flagIDUri struct {
		FlagID uint `uri:"flagID" binding:"required"`
	}
	var flagID flagIDUri
	if err := ctx.ShouldBindUri(&flagID); err != nil {
		ctx.JSONP(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		ctx.Abort()
		return
	}
	flag, ok, msg := db.InitFlagRepo(db.DB.WithContext(ctx)).GetByID(flagID.FlagID, "all")
	if !ok {
		ctx.JSONP(http.StatusOK, gin.H{"msg": msg, "data": nil})
		ctx.Abort()
		return
	}
	ctx.Set("Flag", flag)
	ctx.Next()
}

func GetFlag(ctx *gin.Context) model.Flag {
	if flag, ok := ctx.Get("Flag"); !ok {
		return model.Flag{}
	} else {
		return flag.(model.Flag)
	}
}

// SetContainer 保存 model.Container 至上下文
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
	container, ok, msg := db.InitContainerRepo(db.DB.WithContext(ctx)).GetByID(containerID.ContainerID, "all")
	if !ok {
		ctx.JSONP(http.StatusOK, gin.H{"msg": msg, "data": nil})
		ctx.Abort()
		return
	}
	ctx.Set("Container", container)
	ctx.Next()
}

// GetContainer 从上下文中获取 model.Container
func GetContainer(ctx *gin.Context) model.Container {
	if container, ok := ctx.Get("Container"); !ok {
		return model.Container{}
	} else {
		return container.(model.Container)
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
	notice, ok, msg := db.InitNoticeRepo(db.DB.WithContext(ctx)).GetByID(noticeID.NoticeID, "all")
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
