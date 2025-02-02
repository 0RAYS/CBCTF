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

	// 公共
	router.POST("/register", Register)
	router.POST("/login", Login)
	router.POST("/admin/login", AdminLogin)
	router.GET("/avatar/:avatarID", middleware.SetAvatarID, DownloadAvatar)

	// 鉴权
	auth := router.Group("", middleware.CheckLogin)

	// 用户
	user := auth.Group("/user", middleware.CheckRole("user"))
	{
		user.GET("/info", GetUser)
		user.POST("/password", ChangePassword)
		user.POST("/update", UpdateUser)
		user.POST("/delete", DeleteUser)
		user.POST("/avatar", UploadAvatar(model.User{}))
	}

	// 比赛
	auth.GET("/contest/list", middleware.CheckRole("user"), GetContests)
	contest := auth.Group("/contest/:contestID", middleware.CheckRole("user"), middleware.SetContestID)
	{
		contest.GET("/info", GetContest)
		contest.GET("/team/info", GetTeam)
		contest.GET("/team/captcha", GetTeamCaptcha)
		contest.GET("/team/list", GetTeammates)
		contest.POST("/team/join", middleware.CheckVerified, JoinTeam)
		contest.POST("/team/create", middleware.CheckVerified, CreateTeam)
		contest.POST("/team/update", middleware.CheckVerified, middleware.CheckCaptain, UpdateTeam)
		contest.POST("/team/avatar", middleware.CheckVerified, middleware.CheckCaptain, UploadAvatar(model.Team{}))
		contest.POST("/team/delete", middleware.CheckVerified, middleware.CheckCaptain, DeleteTeam)
		contest.POST("/team/kick", middleware.CheckVerified, middleware.CheckCaptain, KickMember)
		contest.POST("/team/leave", middleware.CheckVerified, LeaveTeam)

		// 比赛题目
		contest.GET("/challenge/list", middleware.CheckVerified, middleware.CheckBanned, GetUsages)
		contestChallenge := contest.Group("/challenge/:challengeID", middleware.CheckVerified, middleware.CheckBanned, middleware.SetChallengeID)
		{
			contestChallenge.POST("/init", InitChallenge)
			contestChallenge.GET("/attachment", GetAttachment)
			contestChallenge.GET("/remote", GetContainer)
			contestChallenge.POST("/start", StartContainer)
			contestChallenge.POST("/increase", IncreaseDuration)
			contestChallenge.POST("/stop", StopContainer)
			contestChallenge.POST("/submit", SubmitFlag)
		}
	}

	// 管理员
	admin := auth.Group("/admin", middleware.CheckRole("admin"))
	{
		// 管理员
		admin.GET("/info", GetAdmin)
		admin.POST("/password", AdminChangePassword)
		admin.POST("/update", UpdateAdmin)
		admin.POST("/avatar", UploadAvatar(model.Admin{}))

		// 系统管理
		adminSystem := admin.Group("/system")
		{
			adminSystem.GET("/status", SystemStatus)
			adminSystem.GET("/config", SystemConfig)
			adminSystem.POST("/update", SystemUpdate)
		}

		// 用户管理
		admin.GET("/user/list", GetUsers)
		admin.POST("/user/create", CreateUser)
		adminUser := admin.Group("/user/:userID", middleware.SetUserID)
		{
			adminUser.GET("/info", GetUser)
			adminUser.POST("/update", UpdateUser)
			adminUser.POST("/delete", DeleteUser)
			adminUser.POST("/avatar", UploadAvatar(model.User{}))
		}

		// 比赛管理
		admin.GET("/contest/list", GetContests)
		admin.POST("/contest/create", CreateContest)
		adminContest := admin.Group("/contest/:contestID", middleware.SetContestID)
		{
			adminContest.GET("/info", GetContest)
			adminContest.GET("/captcha", GetContestCaptcha)
			adminContest.POST("/update", UpdateContest)
			adminContest.POST("/delete", DeleteContest)
			adminContest.POST("/avatar", UploadAvatar(model.Contest{}))
			adminContest.GET("/submissions", GetSubmissions)

			// 比赛队伍管理
			adminContest.GET("/team/list", GetTeams)
			adminContestTeam := adminContest.Group("/team/:teamID", middleware.SetTeamID)
			{
				adminContestTeam.GET("/info", GetTeam)
				adminContestTeam.GET("/captcha", GetTeamCaptcha)
				adminContestTeam.GET("/list", GetTeamUsers)
				adminContestTeam.POST("/update", UpdateTeam)
				adminContestTeam.POST("/delete", DeleteTeam)
				adminContestTeam.POST("/kick", KickMember)
				adminContestTeam.POST("/avatar", UploadAvatar(model.Team{}))
			}

			// 比赛题目管理
			adminContest.GET("/challenge/list", GetUsages)
			adminContest.POST("/challenge/add", AddUsage)
			adminContestChallenge := adminContest.Group("/challenge/:challengeID", middleware.SetChallengeID)
			{
				adminContestChallenge.POST("/update", UpdateUsage)
				adminContestChallenge.POST("/remove", RemoveUsage)
			}
		}

		// 题库管理
		admin.GET("/challenge/list", GetChallenges)
		admin.GET("/challenge/categories", GetCategories)
		admin.POST("/challenge/create", CreateChallenge)
		adminChallenge := admin.Group("/challenge/:challengeID", middleware.SetChallengeID)
		{
			adminChallenge.GET("/info", GetChallenge)
			adminChallenge.GET("/files", GetChallengeFiles)
			adminChallenge.GET("/download", DownloadChallenge)
			adminChallenge.POST("/update", UpdateChallenge)
			adminChallenge.POST("/delete", DeleteChallenge)
			adminChallenge.POST("/upload", UploadChallenge)
		}

		// 头像管理
		admin.GET("/avatar/list", GetAvatars)
		admin.POST("/avatar/delete", DeleteAvatar)
		adminAvatar := admin.Group("/avatar/:avatarID", middleware.SetAvatarID)
		{
			adminAvatar.POST("/delete", DeleteAvatar)
		}
	}

	return router
}
