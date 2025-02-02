package router

import (
	"CBCTF/internal/config"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
	"github.com/gin-gonic/gin"
)

func Init() *gin.Engine {

	gin.SetMode(config.Env.Gin.Mode)
	router := gin.New()
	router.MaxMultipartMemory = int64(config.Env.Gin.Upload.Max << 20)
	router.Use(middleware.Logger(), gin.Recovery(), middleware.Trace, middleware.Cors, middleware.I18n(), middleware.AccessLog, middleware.RateLimit())

	router.POST("/register", Register)
	router.POST("/login", Login)
	router.POST("/admin/login", AdminLogin)
	router.GET("/avatar/:avatarID", middleware.SetAvatarID, DownloadAvatar)
	auth := router.Group("", middleware.CheckLogin)

	user := auth.Group("/user", middleware.CheckRole("user"))
	user.GET("/info", GetUser)
	user.POST("/password", ChangePassword)
	user.POST("/update", UpdateUser)
	user.POST("/delete", DeleteUser)
	user.POST("/avatar", UploadAvatar(model.User{}))

	contest := auth.Group("/contest", middleware.CheckRole("user"))
	contest.GET("/list", GetContests)
	contest.GET("/:contestID/info", middleware.SetContestID, GetContest)
	contest.GET("/:contestID/team/info", middleware.SetContestID, GetTeam)
	contest.GET("/:contestID/team/captcha", middleware.SetContestID, GetTeamCaptcha)
	contest.GET("/:contestID/team/list", middleware.SetContestID, GetTeammates)
	contest.POST("/:contestID/team/join", middleware.CheckVerified, middleware.SetContestID, JoinTeam)
	contest.POST("/:contestID/team/create", middleware.CheckVerified, middleware.SetContestID, CreateTeam)
	contest.POST("/:contestID/team/update", middleware.CheckVerified, middleware.SetContestID, middleware.CheckCaptain, UpdateTeam)
	contest.POST("/:contestID/team/avatar", middleware.CheckVerified, middleware.SetContestID, middleware.CheckCaptain, UploadAvatar(model.Team{}))
	contest.POST("/:contestID/team/delete", middleware.CheckVerified, middleware.SetContestID, middleware.CheckCaptain, DeleteTeam)
	contest.POST("/:contestID/team/kick", middleware.CheckVerified, middleware.SetContestID, middleware.CheckCaptain, KickMember)
	contest.POST("/:contestID/team/leave", middleware.CheckVerified, middleware.SetContestID, LeaveTeam)

	contestChallenge := contest.Group("/:contestID/challenge", middleware.CheckVerified, middleware.SetContestID, middleware.CheckBanned)
	contestChallenge.GET("/list", GetUsages)
	contestChallenge.POST("/:challengeID/init", middleware.SetChallengeID, InitChallenge)
	contestChallenge.GET("/:challengeID/attachment", middleware.SetChallengeID, GetAttachment)
	contestChallenge.GET("/:challengeID/remote", middleware.SetChallengeID, GetContainer)
	contestChallenge.POST("/:challengeID/start", middleware.SetChallengeID, StartContainer)
	contestChallenge.POST("/:challengeID/increase", middleware.SetChallengeID, IncreaseDuration)
	contestChallenge.POST("/:challengeID/stop", middleware.SetChallengeID, StopContainer)
	contestChallenge.POST("/:challengeID/submit", middleware.SetChallengeID, SubmitFlag)

	admin := auth.Group("/admin", middleware.CheckRole("admin"))
	admin.GET("/info", GetAdmin)
	admin.POST("/password", AdminChangePassword)
	admin.POST("/update", UpdateAdmin)
	admin.POST("/avatar", UploadAvatar(model.Admin{}))

	adminSystem := admin.Group("/system")
	adminSystem.GET("/status", SystemStatus)
	adminSystem.GET("/config", SystemConfig)
	adminSystem.POST("/update", SystemUpdate)

	adminUser := admin.Group("/user")
	adminUser.GET("/list", GetUsers)
	adminUser.POST("/create", CreateUser)
	adminUser.GET("/:userID/info", middleware.SetUserID, GetUser)
	adminUser.POST("/:userID/update", middleware.SetUserID, UpdateUser)
	adminUser.POST("/:userID/delete", middleware.SetUserID, DeleteUser)
	adminUser.POST("/:userID/avatar", middleware.SetUserID, UploadAvatar(model.User{}))

	adminContest := admin.Group("/contest")
	adminContest.GET("/list", GetContests)
	adminContest.POST("/create", CreateContest)
	adminContest.GET("/:contestID/info", middleware.SetContestID, GetContest)
	adminContest.GET("/:contestID/captcha", middleware.SetContestID, GetContestCaptcha)
	adminContest.POST("/:contestID/update", middleware.SetContestID, UpdateContest)
	adminContest.POST("/:contestID/delete", middleware.SetContestID, DeleteContest)
	adminContest.POST("/:contestID/avatar", middleware.SetContestID, UploadAvatar(model.Contest{}))

	adminContestTeam := adminContest.Group("/:contestID/team", middleware.SetContestID)
	adminContestTeam.GET("/list", GetTeams)
	adminContestTeam.GET("/:teamID/info", middleware.SetTeamID, GetTeam)
	adminContestTeam.GET("/:teamID/captcha", middleware.SetTeamID, GetTeamCaptcha)
	adminContestTeam.GET("/:teamID/list", middleware.SetTeamID, GetTeamUsers)
	adminContestTeam.POST("/:teamID/update", middleware.SetTeamID, UpdateTeam)
	adminContestTeam.POST("/:teamID/delete", middleware.SetTeamID, DeleteTeam)
	adminContestTeam.POST("/:teamID/kick", middleware.SetTeamID, KickMember)
	adminContestTeam.POST("/:teamID/avatar", middleware.SetTeamID, UploadAvatar(model.Team{}))

	adminContestChallenge := adminContest.Group("/:contestID/challenge", middleware.SetContestID)
	adminContestChallenge.GET("/list", GetUsages)
	adminContestChallenge.POST("/add", AddUsage)
	adminContestChallenge.POST("/:challengeID/update", middleware.SetChallengeID, UpdateUsage)
	adminContestChallenge.POST("/:challengeID/remove", middleware.SetChallengeID, RemoveUsage)

	adminChallenge := admin.Group("/challenge")
	adminChallenge.GET("/list", GetChallenges)
	adminChallenge.GET("/categories", GetCategories)
	adminChallenge.GET("/:challengeID/info", middleware.SetChallengeID, GetChallenge)
	adminChallenge.GET("/:challengeID/files", middleware.SetChallengeID, GetChallengeFiles)
	adminChallenge.GET("/:challengeID/download", middleware.SetChallengeID, DownloadChallenge)
	adminChallenge.POST("/create", CreateChallenge)
	adminChallenge.POST("/:challengeID/update", middleware.SetChallengeID, UpdateChallenge)
	adminChallenge.POST("/:challengeID/delete", middleware.SetChallengeID, DeleteChallenge)
	adminChallenge.POST("/:challengeID/upload", middleware.SetChallengeID, UploadChallenge)

	adminFile := admin.Group("/avatar")
	adminFile.GET("/list", GetAvatars)
	adminFile.POST("/delete", DeleteAvatar)
	adminFile.POST("/:avatarID/delete", middleware.SetAvatarID, DeleteAvatar)

	return router

}
