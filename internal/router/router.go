package router

import (
	"CBCTF/internal/config"
	"CBCTF/internal/log"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
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
	router.Use(middleware.Logger(), gin.Recovery(), middleware.Trace, middleware.Cors, middleware.I18n(), middleware.AccessLog, middleware.RateLimit())

	// 公共
	router.POST("/register", Register)
	router.POST("/login", Login)
	router.POST("/admin/login", AdminLogin)
	router.GET("/verify", Verify)
	router.GET("/avatar/:fileID", middleware.SetFile(model.Avatar), DownloadFile)

	router.GET("/contests", GetContests)

	// 鉴权
	auth := router.Group("", middleware.CheckLogin)

	// 用户
	user := auth.Group("/me", middleware.CheckRole("user"))
	{
		user.GET("", GetUser)
		user.PUT("/password", ChangePassword)
		user.PUT("", UpdateUser)
		user.DELETE("", DeleteUser)
		user.PUT("/avatar", UploadAvatar("self-user"))
		user.POST("/activate", Activate)
	}

	// 比赛
	contest := auth.Group("/contests/:contestID", middleware.CheckRole("user"), middleware.SetContest)
	{
		contest.GET("", GetContest)
		contest.GET("/rank", GetRank)
		contest.GET("/rank/detail", GetRankDetail)
		contest.POST("/teams/join", middleware.CheckRunning, middleware.CheckVerified, JoinTeam)
		contest.POST("/teams/create", middleware.CheckRunning, middleware.CheckVerified, CreateTeam)

		contestTeam := contest.Group("/teams", middleware.SetTeamByUser)
		{
			contestTeam.GET("/me", GetTeam)
			contestTeam.GET("/me/captcha", GetTeamCaptcha)
			contestTeam.GET("/me/users", GetTeammates)
			contestTeam.PUT("/me/captcha", UpdateCaptcha)
			contestTeam.PUT("/me", middleware.CheckRunning, middleware.CheckVerified, middleware.CheckCaptain, UpdateTeam)
			contestTeam.PUT("/me/avatar", middleware.CheckRunning, middleware.CheckVerified, middleware.CheckCaptain, UploadAvatar("team"))
			contestTeam.DELETE("/me", middleware.CheckRunning, middleware.CheckVerified, middleware.CheckCaptain, DeleteTeam)
			contestTeam.POST("/me/kick", middleware.CheckRunning, middleware.CheckVerified, middleware.CheckCaptain, KickMember)
			contestTeam.POST("/me/leave", middleware.CheckRunning, middleware.CheckVerified, LeaveTeam)
		}

		// 比赛公告
		{
			contest.GET("/notices", GetNotices)
			contest.GET("/notices/:noticeID", middleware.SetNotice, GetNotice)
		}

		// 比赛题目
		contest.GET("/challenges", middleware.CheckVerified, middleware.SetTeamByUser, middleware.CheckBanned, GetUsages)
		contestChallenge := contest.Group("/challenges/:challengeID", middleware.CheckVerified, middleware.SetTeamByUser, middleware.CheckBanned, middleware.SetChallenge)
		{
			contestChallenge.GET("", ChallengeStatus)
			contestChallenge.POST("/init", middleware.CheckRunning, InitChallenge(false))
			contestChallenge.GET("/attachment", GetAttachment)
			contestChallenge.POST("/reset", middleware.CheckRunning, InitChallenge(true))
			contestChallenge.POST("/start", StartContainer)
			contestChallenge.POST("/increase", middleware.CheckRunning, IncreaseDuration)
			contestChallenge.POST("/stop", StopContainer)
			contestChallenge.POST("/submit", middleware.CheckRunning, SubmitFlag)
		}

		// WriteUp
		contestWriteUp := contest.Group("/writeups", middleware.CheckVerified, middleware.SetTeamByUser, middleware.CheckBanned, middleware.CheckRunning)
		{
			contestWriteUp.POST("", UploadWriteUp)
			contestWriteUp.GET("", GetWriteUPs)
		}
	}

	// 管理员
	admin := auth.Group("/admin", middleware.CheckRole("admin"))
	{
		// 管理员
		admin.GET("/me", GetAdmin)
		admin.GET("/admins", GetAdmins)
		admin.PUT("/me/password", AdminChangePassword)
		admin.PUT("/me/update", UpdateAdmin)
		admin.PUT("/me/avatar", UploadAvatar("self-admin"))
		admin.POST("/admins", CreateAdmin)

		// 系统管理
		adminSystem := admin.Group("/system")
		{
			adminSystem.GET("/status", SystemStatus)
			adminSystem.GET("/config", SystemConfig)
			//adminSystem.POST("/update", SystemUpdate)
		}

		// 用户管理
		admin.GET("/users", GetUsers)
		admin.POST("/users", CreateUser)
		adminUser := admin.Group("/users/:userID", middleware.SetUser)
		{
			adminUser.GET("", GetUser)
			adminUser.PUT("", UpdateUser)
			adminUser.DELETE("", DeleteUser)
			adminUser.PUT("/avatar", UploadAvatar("user"))
		}

		// 比赛管理
		admin.GET("/contests", GetContests)
		admin.POST("/contests", CreateContest)
		adminContest := admin.Group("/contests/:contestID", middleware.SetContest)
		{
			adminContest.GET("", GetContest)
			adminContest.GET("/captcha", GetContestCaptcha)
			adminContest.PUT("", UpdateContest)
			adminContest.DELETE("", DeleteContest)
			adminContest.PUT("/avatar", UploadAvatar("contest"))
			adminContest.GET("/submissions", GetSubmissions)
			adminContest.GET("/rank", GetRank)
			adminContest.GET("/rank/detail", GetRankDetail)

			// 比赛队伍管理
			adminContest.GET("/teams", GetTeams)
			adminContestTeam := adminContest.Group("/teams/:teamID", middleware.SetTeamByURI)
			{
				adminContestTeam.GET("", GetTeam)
				adminContestTeam.GET("/captcha", GetTeamCaptcha)
				adminContestTeam.GET("/users", GetTeammates)
				adminContestTeam.PUT("", UpdateTeam)
				adminContestTeam.DELETE("", DeleteTeam)
				adminContestTeam.POST("/kick", KickMember)
				adminContestTeam.POST("/avatar", UploadAvatar("team"))

				// 比赛队伍提交管理
				adminContestTeam.GET("/submissions", GetTeamSubmissions)

				// 比赛队伍容器管理
				adminContestTeam.GET("/containers", GetContainers)
				adminContainer := adminContestTeam.Group("/containers/:containerID", middleware.SetContainer)
				{
					adminContainer.GET("", GetContainer)
					adminContainer.POST("/stop", StopContainer)

					adminTraffic := adminContainer.Group("/traffic")
					adminTraffic.GET("/download", DownloadTraffic)
					adminTraffic.POST("/load", LoadTraffic)
					adminTraffic.GET("", GetTraffics)
				}

				// WriteUp
				adminContestTeam.GET("/writeups", GetWriteUPs)
				adminContestTeam.GET("/writeups/:fileID", middleware.SetFile(model.WriteUP), DownloadFile)
			}

			// 比赛公告管理
			adminContest.GET("/notices", GetNotices)
			adminContest.POST("/notices", CreateNotice)
			adminContestNotice := adminContest.Group("/notices/:noticeID", middleware.SetNotice)
			{
				adminContestNotice.GET("", GetNotice)
				adminContestNotice.PUT("", UpdateNotice)
				adminContestNotice.DELETE("", DeleteNotice)
			}

			// 比赛题目管理
			adminContest.GET("/challenges", GetUsages)
			adminContest.POST("/challenges", AddUsage)
			adminContestChallenge := adminContest.Group("/challenges/:challengeID", middleware.SetChallenge)
			{
				adminContestChallenge.PUT("", UpdateUsage)
				adminContestChallenge.DELETE("", RemoveUsage)
			}
		}

		// 题库管理
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

		// 头像管理
		admin.GET("/avatars", GetAvatars)
		admin.DELETE("/avatars", DeleteFile)
		adminAvatar := admin.Group("/avatars/:fileID", middleware.SetFile(model.Avatar))
		{
			adminAvatar.DELETE("/delete", DeleteFile)
		}
	}

	return router
}
