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
	router.GET("/avatar/:fileID", middleware.SetFile, DownloadFile)

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
		user.POST("/activate", Activate)
	}

	// 比赛
	auth.GET("/contest/list", middleware.CheckRole("user"), GetContests)
	contest := auth.Group("/contest/:contestID", middleware.CheckRole("user"), middleware.SetContest)
	{
		contest.GET("/info", GetContest)
		contest.GET("/rank", GetRank)
		contest.GET("/rank/detail", GetRankDetail)
		contest.POST("/join", middleware.CheckVerified, JoinTeam)
		contest.POST("/create", middleware.CheckVerified, CreateTeam)

		contestTeam := contest.Group("/team", middleware.SetTeamByUser)
		{
			contestTeam.GET("/info", GetTeam)
			contestTeam.GET("/captcha", GetTeamCaptcha)
			contestTeam.GET("/list", GetTeammates)
			contestTeam.POST("/captcha", UpdateCaptcha)
			contestTeam.POST("/update", middleware.CheckVerified, middleware.CheckCaptain, UpdateTeam)
			contestTeam.POST("/avatar", middleware.CheckVerified, middleware.CheckCaptain, UploadAvatar(model.Team{}))
			contestTeam.POST("/delete", middleware.CheckVerified, middleware.CheckCaptain, DeleteTeam)
			contestTeam.POST("/kick", middleware.CheckVerified, middleware.CheckCaptain, KickMember)
			contestTeam.POST("/leave", middleware.CheckVerified, LeaveTeam)
		}

		// 比赛公告
		{
			contest.GET("/notice/list", GetNotices)
			contest.GET("/notice/:noticeID/info", middleware.SetNotice, GetNotice)
		}

		// 比赛题目
		contest.GET("/challenge/list", middleware.CheckVerified, middleware.SetTeamByUser, middleware.CheckBanned, middleware.CheckRunning, GetUsages)
		contestChallenge := contest.Group("/challenge/:challengeID", middleware.CheckVerified, middleware.SetTeamByUser, middleware.CheckBanned, middleware.CheckRunning, middleware.SetChallenge)
		{
			contestChallenge.GET("/status", ChallengeStatus)
			contestChallenge.POST("/init", InitChallenge(false))
			contestChallenge.GET("/files", GetChallengeFiles)
			contestChallenge.GET("/attachment", GetAttachment)
			contestChallenge.GET("/remote", GetContainer(false))
			contestChallenge.POST("/reset", InitChallenge(true))
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
		admin.GET("/list", GetAdmins)
		admin.POST("/password", AdminChangePassword)
		admin.POST("/update", UpdateAdmin)
		admin.POST("/avatar", UploadAvatar(model.Admin{}))
		admin.POST("/create", CreateAdmin)

		// 系统管理
		adminSystem := admin.Group("/system")
		{
			adminSystem.GET("/status", SystemStatus)
			adminSystem.GET("/config", SystemConfig)
			//adminSystem.POST("/update", SystemUpdate)
		}

		// 用户管理
		admin.GET("/user/list", GetUsers)
		admin.POST("/user/create", CreateUser)
		adminUser := admin.Group("/user/:userID", middleware.SetUser)
		{
			adminUser.GET("/info", GetUser)
			adminUser.POST("/update", UpdateUser)
			adminUser.POST("/delete", DeleteUser)
			adminUser.POST("/avatar", UploadAvatar(model.User{}))
		}

		// 比赛管理
		admin.GET("/contest/list", GetContests)
		admin.POST("/contest/create", CreateContest)
		adminContest := admin.Group("/contest/:contestID", middleware.SetContest)
		{
			adminContest.GET("/info", GetContest)
			adminContest.GET("/captcha", GetContestCaptcha)
			adminContest.POST("/update", UpdateContest)
			adminContest.POST("/delete", DeleteContest)
			adminContest.POST("/avatar", UploadAvatar(model.Contest{}))
			adminContest.GET("/submissions", GetSubmissions)
			adminContest.GET("/rank", GetRank)
			adminContest.GET("/rank/detail", GetRankDetail)

			// 比赛队伍管理
			adminContest.GET("/team/list", GetTeams)
			adminContestTeam := adminContest.Group("/team/:teamID", middleware.SetTeamByURI)
			{
				adminContestTeam.GET("/info", GetTeam)
				adminContestTeam.GET("/captcha", GetTeamCaptcha)
				adminContestTeam.GET("/list", GetTeammates)
				adminContestTeam.POST("/update", UpdateTeam)
				adminContestTeam.POST("/delete", DeleteTeam)
				adminContestTeam.POST("/kick", KickMember)
				adminContestTeam.POST("/avatar", UploadAvatar(model.Team{}))

				// 比赛队伍提交管理
				adminContestTeam.GET("/submission/list", GetTeamSubmissions)

				// 比赛队伍容器管理
				adminContestTeam.GET("/container/list", GetContainers)
				adminContainer := adminContestTeam.Group("/container/:containerID", middleware.SetContainer)
				{
					adminContainer.GET("/info", GetContainer(true))
					adminContainer.POST("/stop", StopContainer)

					adminTraffic := adminContainer.Group("/traffic")
					adminTraffic.GET("/download", DownloadTraffic)
					adminTraffic.POST("/load", LoadTraffic)
					adminTraffic.GET("/show", GetTraffics)
				}
			}

			// 比赛公告管理
			adminContest.GET("/notice/list", GetNotices)
			adminContest.POST("/notice/create", CreateNotice)
			adminContestNotice := adminContest.Group("/notice/:noticeID", middleware.SetNotice)
			{
				adminContestNotice.GET("/info", GetNotice)
				adminContestNotice.POST("/update", UpdateNotice)
				adminContestNotice.POST("/delete", DeleteNotice)
			}

			// 比赛题目管理
			adminContest.GET("/challenge/list", GetUsages)
			adminContest.POST("/challenge/add", AddUsage)
			adminContestChallenge := adminContest.Group("/challenge/:challengeID", middleware.SetChallenge)
			{
				adminContestChallenge.POST("/update", UpdateUsage)
				adminContestChallenge.POST("/remove", RemoveUsage)
			}
		}

		// 题库管理
		admin.GET("/challenge/list", GetChallenges)
		admin.GET("/challenge/categories", GetCategories)
		admin.POST("/challenge/create", CreateChallenge)
		adminChallenge := admin.Group("/challenge/:challengeID", middleware.SetChallenge)
		{
			adminChallenge.GET("/info", GetChallenge)
			adminChallenge.GET("/files", GetChallengeFiles)
			adminChallenge.GET("/download", DownloadChallenge)
			adminChallenge.POST("/update", UpdateChallenge)
			adminChallenge.POST("/delete", DeleteChallenge)
			adminChallenge.POST("/upload", UploadChallenge)
		}

		// 头像管理
		admin.GET("/avatar/list", GetFiles("image"))
		admin.POST("/avatar/delete", DeleteFile)
		adminAvatar := admin.Group("/avatar/:fileID", middleware.SetFile)
		{
			adminAvatar.POST("/delete", DeleteFile)
		}
	}

	return router
}
