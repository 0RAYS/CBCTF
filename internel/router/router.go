package router

import (
	"CBCTF/internel/config"
	"CBCTF/internel/log"
	"CBCTF/internel/middleware"
	"CBCTF/internel/model"
	"github.com/gin-gonic/gin"
)

func Init() *gin.Engine {
	gin.SetMode(config.Env.Gin.Mode)
	router := gin.New()
	log.Logger.Infof("Trust proxies: %v", config.Env.Gin.Proxies)
	if err := router.SetTrustedProxies(config.Env.Gin.Proxies); err != nil {
		log.Logger.Warningf("Set trusted proxies failed: %v", err)
	}
	router.MaxMultipartMemory = int64(config.Env.Gin.Upload.Max << 20)
	router.Use(
		middleware.Logger(), gin.Recovery(), middleware.SetTrace, middleware.Cors,
		middleware.I18n(), middleware.AccessLog, middleware.RateLimit(), middleware.SetMagic,
	)

	{
		router.POST("/register", Register)
		router.POST("/login", Login)
		router.POST("/admin/login", AdminLogin)
		router.POST("/verify", VerifyEmail)

		router.GET("/stats", HomePage)
		router.GET("/contests", GetContests)
	}

	auth := router.Group("", middleware.CheckAuth)

	user := auth.Group("/me", middleware.CheckRole("user"))
	{
		user.GET("", GetUser)
		user.PUT("/password", ChangePwd)
		user.PUT("", UpdateUser)
		user.DELETE("", DeleteUser)
		user.PUT("/avatar", UploadFile("self-user", "avatar"))
		user.POST("/activate", ActivateEmail)
	}

	contest := auth.Group("/contests/:contestID", middleware.CheckRole("user"), middleware.SetContest)
	{
		contest.GET("", GetContest)
		contest.GET("/rank", GetTeamRanking)
		contest.POST("/teams/join", middleware.CheckRunning, middleware.CheckVerified, JoinTeam)
		contest.POST("/teams/create", middleware.CheckRunning, middleware.CheckVerified, CreateTeam)

		contestTeam := contest.Group("/teams", middleware.SetTeamByUser)
		{
			contestTeam.GET("/me", GetTeam)
			contestTeam.GET("/me/captcha", GetTeamCaptcha)
			contestTeam.GET("/me/users", GetTeammates)
			contestTeam.PUT("/me/captcha", UpdateCaptcha)
			contestTeam.PUT("/me", middleware.CheckRunning, middleware.CheckVerified, middleware.CheckCaptain, UpdateTeam)
			contestTeam.PUT("/me/avatar", middleware.CheckRunning, middleware.CheckVerified, middleware.CheckCaptain, UploadFile("team", "avatar"))
			contestTeam.DELETE("/me", middleware.CheckRunning, middleware.CheckVerified, middleware.CheckCaptain, DeleteTeam)
			contestTeam.POST("/me/kick", middleware.CheckRunning, middleware.CheckVerified, middleware.CheckCaptain, KickMember)
			contestTeam.POST("/me/leave", middleware.CheckRunning, LeaveTeam)
		}

		// 比赛公告
		{
			contest.GET("/notices", GetNotices)
			contest.GET("/notices/:noticeID", middleware.SetNotice, GetNotice)
		}

		// 比赛题目
		contest.GET("/challenges", middleware.CheckVerified, middleware.SetTeamByUser, middleware.CheckBanned, GetUsages)
		contestChallenge := contest.Group(
			"/challenges/:challengeID",
			middleware.CheckVerified, middleware.SetTeamByUser, middleware.CheckBanned, middleware.SetUsage,
		)
		{
			contestChallenge.GET("", GetUsageStatus)
			contestChallenge.POST("/init", middleware.CheckRunning, middleware.CheckGenerated, InitUsage(false))
			contestChallenge.GET("/attachment", GetAttachment)
			contestChallenge.POST("/reset", middleware.CheckRunning, middleware.CheckGenerated, InitUsage(true))
			//TODO container op
		}

		//TODO writeup CURD
	}

	admin := auth.Group("/admin", middleware.CheckRole("admin"))
	{
		admin.GET("/me", GetAdmin)
		admin.PUT("/me/password", AdminChangePassword)
		admin.PUT("/me/update", UpdateAdmin)
		admin.PUT("/me/avatar", UploadFile("admin", "avatar"))
		admin.POST("/admins", CreateAdmin)

		// 系统管理
		adminSystem := admin.Group("/system")
		{
			adminSystem.GET("/status", SystemStatus)
			adminSystem.GET("/config", SystemConfig)
		}

		admin.GET("/users", GetUsers)
		admin.POST("/users", CreateUser)
		adminUser := admin.Group("/users/:userID", middleware.SetUser)
		{
			adminUser.GET("", GetUser)
			adminUser.PUT("", UpdateUser)
			adminUser.DELETE("", DeleteUser)
			adminUser.PUT("/avatar", UploadFile("user", "avatar"))
		}

		admin.GET("/contests", GetContests)
		admin.POST("/contests", CreateContest)
		adminContest := admin.Group("/contests/:contestID", middleware.SetContest)
		{
			adminContest.GET("", GetContest)
			adminContest.PUT("", UpdateContest)
			adminContest.DELETE("", DeleteContest)
			adminContest.PUT("/avatar", UploadFile("contest", "avatar"))
			adminContest.GET("/submissions", GetSubmissions(false))
			adminContest.GET("/rank", GetRank)
		}

		adminContest.GET("/teams", GetTeams)
		adminContestTeam := adminContest.Group("/teams/:teamID", middleware.SetTeamByURI)
		{
			adminContestTeam.GET("", GetTeam)
			adminContestTeam.GET("/captcha", GetTeamCaptcha)
			adminContestTeam.GET("/users", GetTeammates)
			adminContestTeam.PUT("", UpdateTeam)
			adminContestTeam.DELETE("", DeleteTeam)
			adminContestTeam.POST("/kick", KickMember)
			adminContestTeam.POST("/avatar", UploadFile("team", "avatar"))

			adminContestTeam.GET("/submissions", GetSubmissions(true))

			adminContestTeam.GET("/containers", GetContainers)
			adminContainer := adminContestTeam.Group("/containers/:containerID", middleware.SetContainer)
			{
				adminContainer.GET("", GetContainer)
				//adminContainer.POST("/stop", StopContainer)

				adminTraffic := adminContainer.Group("/traffic")
				adminTraffic.GET("/download", DownloadTraffic)
				adminTraffic.POST("/load", LoadTraffic)
				adminTraffic.GET("", GetTraffics)
			}

			adminContestTeam.GET("/writeups", GetWriteUPs)
			adminContestTeam.GET("/writeups/:fileID", middleware.SetFile(model.WriteUP), DownloadFile)
		}

		adminContest.GET("/notices", GetNotices)
		adminContest.POST("/notices", CreateNotice)
		adminContestNotice := adminContest.Group("/notices/:noticeID", middleware.SetNotice)
		{
			adminContestNotice.GET("", GetNotice)
			adminContestNotice.PUT("", UpdateNotice)
			adminContestNotice.DELETE("", DeleteNotice)
		}

		adminContest.GET("/challenges", GetUsages)
		adminContest.POST("/challenges", AddUsage)
		adminContestChallenge := adminContest.Group("/challenges/:challengeID", middleware.SetUsage)
		{
			adminContestChallenge.PUT("", UpdateUsage)
			adminContestChallenge.DELETE("", RemoveUsage)

			//TODO flag CURD
		}

	}

	return router
}
