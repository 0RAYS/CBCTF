package router

import (
	"CBCTF/internel/config"
	"CBCTF/internel/log"
	"CBCTF/internel/middleware"
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

		}
	}

	return router
}
