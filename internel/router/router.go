package router

import (
	"CBCTF/internel/config"
	"CBCTF/internel/log"
	"CBCTF/internel/middleware"
	"CBCTF/internel/model"
	"errors"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	_ "net/http/pprof"
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
		gin.Recovery(), middleware.Cors, middleware.Logger, middleware.Prometheus, middleware.SetTrace,
		middleware.SetMagic, middleware.I18n, middleware.AccessLog, middleware.RateLimit, middleware.Events,
	)

	{
		pprof.Register(router)

		prometheus.MustRegister(middleware.HttpRequestsTotal)
		prometheus.MustRegister(middleware.HttpRequestDuration)
		prometheus.MustRegister(middleware.HttpRequestSize)
		prometheus.MustRegister(middleware.HttpResponseSize)
		prometheus.MustRegister(middleware.InFlightRequests)
		var alreadyRegisteredError prometheus.AlreadyRegisteredError
		if err := prometheus.Register(collectors.NewGoCollector()); err != nil {
			if !errors.As(err, &alreadyRegisteredError) {
				log.Logger.Warningf("failed to register GoCollector: %v", err)
			}
		}
		if err := prometheus.Register(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{})); err != nil {
			if !errors.As(err, &alreadyRegisteredError) {
				log.Logger.Warningf("failed to register ProcessCollector: %v", err)
			}
		}
		router.GET("/metrics", gin.WrapH(promhttp.Handler()))
	}

	{
		router.POST("/register", Register)
		router.POST("/login", Login)
		router.POST("/admin/login", AdminLogin)
		router.POST("/verify", VerifyEmail)
		router.GET("/avatars/:fileID", middleware.SetFile(model.Avatar), DownloadFile)

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
		user.POST("/avatar", UploadAvatar("self-user"))
		user.POST("/activate", ActivateEmail)
	}

	contest := auth.Group("/contests/:contestID", middleware.CheckRole("user"), middleware.SetContest)
	{
		contest.GET("", GetContest)
		contest.GET("/rank", GetTeamRanking)
		contest.POST("/teams/join", middleware.ContestIsNotOver, middleware.CheckVerified, JoinTeam)
		contest.POST("/teams/create", middleware.ContestIsNotOver, middleware.CheckVerified, CreateTeam)

		contestTeam := contest.Group("/teams/me", middleware.SetTeamByUser)
		{
			contestTeam.GET("", GetTeam)
			contestTeam.GET("/captcha", GetTeamCaptcha)
			contestTeam.GET("/users", GetTeammates)
			contestTeam.PUT("/captcha", middleware.ContestIsNotOver, middleware.CheckVerified, middleware.CheckCaptain, UpdateCaptcha)
			contestTeam.PUT("", middleware.ContestIsNotOver, middleware.CheckVerified, middleware.CheckCaptain, UpdateTeam)
			contestTeam.POST("/avatar", middleware.ContestIsNotOver, middleware.CheckVerified, middleware.CheckCaptain, UploadAvatar("team"))
			contestTeam.DELETE("", middleware.ContestStatus(model.ContestIsComing), middleware.CheckVerified, middleware.CheckCaptain, DeleteTeam)
			contestTeam.POST("/kick", middleware.ContestStatus(model.ContestIsComing), middleware.CheckVerified, middleware.CheckCaptain, KickMember)
			contestTeam.POST("/leave", middleware.ContestStatus(model.ContestIsComing), LeaveTeam)
		}

		// 比赛公告
		{
			contest.GET("/notices", GetNotices)
			contest.GET("/notices/:noticeID", middleware.SetNotice, GetNotice)
		}

		// 比赛题目
		contest.GET("/challenges", middleware.CheckVerified, middleware.SetTeamByUser, middleware.CheckBanned, middleware.ContestIsNotComing, GetUsages)
		contestChallenge := contest.Group(
			"/challenges/:challengeID",
			middleware.CheckVerified, middleware.SetTeamByUser, middleware.CheckBanned, middleware.ContestIsNotComing, middleware.SetUsage,
		)
		{
			contestChallenge.GET("", GetUsageStatus)
			contestChallenge.POST("/init", middleware.ContestStatus(model.ContestIsRunning), middleware.CheckVerified, middleware.CheckCaptain, middleware.CheckSolved, GenerateTeamUsage(false))
			contestChallenge.GET("/attachment", DownloadAttachment)
			contestChallenge.POST("/reset", middleware.ContestStatus(model.ContestIsRunning), middleware.CheckGenerated, middleware.CheckSolved, GenerateTeamUsage(true))
			contestChallenge.POST("/start", middleware.CheckGenerated, StartVictim)
			contestChallenge.POST("/increase", middleware.ContestStatus(model.ContestIsRunning), middleware.CheckGenerated, IncreaseVictimDuration)
			contestChallenge.POST("/stop", middleware.CheckGenerated, StopVictim)
			contestChallenge.POST("/submit", middleware.ContestStatus(model.ContestIsRunning), middleware.CheckGenerated, middleware.CheckSolved, SubmitFlag)
		}

		// WriteUp
		contestWriteUp := contest.Group(
			"/writeups",
			middleware.CheckVerified, middleware.SetTeamByUser, middleware.CheckBanned, middleware.ContestIsNotComing,
		)
		{
			contestWriteUp.POST("", UploadWriteUp)
			contestWriteUp.GET("", GetWriteUPs)
		}
	}

	admin := auth.Group("/admin", middleware.CheckRole("admin"))
	{
		admin.GET("/me", GetAdmin)
		admin.PUT("/me/password", AdminChangePassword)
		admin.PUT("/me", UpdateAdmin)
		admin.POST("/me/avatar", UploadAvatar("admin"))
		admin.POST("", CreateAdmin)

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
			adminUser.POST("/avatar", UploadAvatar("user"))
		}

		admin.GET("/contests", GetContests)
		admin.POST("/contests", CreateContest)
		adminContest := admin.Group("/contests/:contestID", middleware.SetContest)
		{
			adminContest.GET("", GetContest)
			adminContest.PUT("", UpdateContest)
			adminContest.DELETE("", DeleteContest)
			adminContest.POST("/avatar", UploadAvatar("contest"))
			adminContest.GET("/rank", GetTeamRanking)

			adminContest.GET("/teams", GetTeams)
			adminContestTeam := adminContest.Group("/teams/:teamID", middleware.SetTeamByURI)
			{
				adminContestTeam.GET("", GetTeam)
				adminContestTeam.GET("/users", GetTeammates)
				adminContestTeam.PUT("", UpdateTeam)
				adminContestTeam.DELETE("", DeleteTeam)
				adminContestTeam.POST("/kick", KickMember)
				adminContestTeam.POST("/avatar", UploadAvatar("team"))

				adminContestTeam.GET("/submissions", GetSubmissions)

				adminContestTeam.GET("/victims", GetVictims)
				adminContainer := adminContestTeam.Group("/victims/:victimID", middleware.SetVictim)
				{
					adminContainer.GET("", GetVictim)

					adminTraffic := adminContainer.Group("/traffic")
					adminTraffic.GET("/download", DownloadTraffic)
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
			adminContestUsage := adminContest.Group("/challenges/:challengeID", middleware.SetUsage)
			{
				adminContestUsage.GET("", GetUsage)
				adminContestUsage.PUT("", UpdateUsage)
				adminContestUsage.DELETE("", RemoveUsage)

				//不允许后期创建和删除
				adminContestUsage.GET("/flags", GetFlags)
				adminContestUsageFlag := adminContestUsage.Group("/flags/:flagID", middleware.SetFlag)
				{
					adminContestUsageFlag.GET("", GetFlag)
					adminContestUsageFlag.PUT("", UpdateFlag)
				}
			}
		}

		admin.GET("/challenges", GetChallenges)
		admin.GET("/challenges/categories", GetCategories)
		admin.POST("/challenges", CreateChallenge)
		adminChallenge := admin.Group("/challenges/:challengeID", middleware.SetChallenge)
		{
			adminChallenge.GET("", GetChallenge)
			adminChallenge.GET("/files", GetChallengeFiles)
			adminChallenge.GET("/download", DownloadChallenge)
			adminChallenge.PUT("", UpdateChallenge)
			adminChallenge.DELETE("", DeleteChallenge)
			adminChallenge.POST("/upload", UploadChallenge)
		}

		admin.GET("/avatars", GetAvatars)
		admin.DELETE("/avatars", DeleteAvatars)
	}

	return router
}
