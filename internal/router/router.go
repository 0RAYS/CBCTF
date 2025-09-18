package router

import (
	"CBCTF/frontend"
	"CBCTF/internal/config"
	"CBCTF/internal/log"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
	"CBCTF/internal/websocket"
	wsm "CBCTF/internal/websocket/middleware"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

func Init() *gin.Engine {
	gin.SetMode(config.Env.Gin.Mode)
	router := gin.New()

	log.Logger.Infof("Trust proxies: %s", config.Env.Gin.Proxies)
	if err := router.SetTrustedProxies(config.Env.Gin.Proxies); err != nil {
		log.Logger.Warningf("Set trusted proxies failed: %s", err)
	}

	router.MaxMultipartMemory = int64(config.Env.Gin.Upload.Max << 20)

	router.Use(gin.Recovery(), middleware.Cors())

	{
		// 不可接入其他中间件
		router.GET("/ws", wsm.SetTrace, wsm.SetMagic, wsm.WSAuth, websocket.WS)
	}

	{
		router.GET("/", func(ctx *gin.Context) {
			if strings.HasPrefix(config.Env.Frontend, "http://") || strings.HasPrefix(config.Env.Frontend, "https://") {
				ctx.Redirect(http.StatusFound, fmt.Sprintf("%s/platform", config.Env.Frontend))
			} else {
				ctx.Redirect(http.StatusFound, "/platform")
			}
		})
		router.StaticFS("/platform", http.FS(frontend.SubFS))
	}

	router.Use(
		middleware.SetTrace, middleware.SetMagic, middleware.Logger, middleware.Prometheus, middleware.AccessLog,
		middleware.I18n, middleware.RateLimit("globals", config.Env.Gin.RateLimit.Global, time.Minute), middleware.Events,
	)

	{
		pprof.Register(router)
		RegisterMetricsRouter(router)
	}

	{
		router.POST("/register", middleware.RateLimit("register", 1, time.Minute), Register)
		router.POST("/login", Login)
		router.POST("/admin/login", AdminLogin)

		RegisterOauthRouter()
		router.GET("/oauth", ListOauth)
		oauth := router.Group("/oauth/:oauth", middleware.SetOauthUri)
		{
			oauth.GET("", Oauth)
			oauth.GET("/callback", OauthCallback)
		}

		router.GET("/verify", VerifyEmail)
		router.GET("/assets", DefaultAssets)
		router.GET("/avatars/:fileID", middleware.SetFile(model.AvatarFileType), DownloadFile(model.SkipEventType))

		router.GET("/stats", HomePage)
		router.GET("/contests", GetContests)
	}

	auth := router.Group("", middleware.CheckAuth)

	user := auth.Group("/me", middleware.CheckRole(false))
	{
		user.GET("", GetUser)
		user.PUT("/password", ChangePwd)
		user.PUT("", UpdateUser)
		user.DELETE("", DeleteUser)
		user.POST("/avatar", UploadAvatar("self-user"))
		user.POST("/activate", middleware.RateLimit("activate", 1, time.Minute), ActivateEmail)
	}

	contest := auth.Group("/contests/:contestID", middleware.CheckRole(false), middleware.SetContest)
	{
		contest.GET("", GetContest)
		contest.GET("/rank", GetTeamRanking)
		contest.GET("/scoreboard", GetScoreboard)
		contest.GET("/timeline", GetRankTimeline)
		contest.POST("/teams/join", middleware.ContestIsNotOver, middleware.CheckVerified, JoinTeam)
		contest.POST("/teams/create", middleware.ContestIsNotOver, middleware.CheckVerified, CreateTeam)

		contestTeam := contest.Group("/teams/me", middleware.CheckVerified, middleware.SetTeamByUser)
		{
			contestTeam.GET("", GetTeam)
			contestTeam.GET("/captcha", GetTeamCaptcha)
			contestTeam.GET("/users", GetTeammates)
			contestTeam.PUT("/captcha", middleware.ContestIsNotOver, middleware.CheckCaptain, UpdateCaptcha)
			contestTeam.PUT("", middleware.ContestIsNotOver, middleware.CheckCaptain, UpdateTeam)
			contestTeam.POST("/avatar", middleware.ContestIsNotOver, middleware.CheckCaptain, UploadAvatar("team"))
			contestTeam.DELETE("", middleware.ContestIsComing, middleware.CheckCaptain, DeleteTeam)
			contestTeam.POST("/kick", middleware.ContestIsComing, middleware.CheckCaptain, KickMember)
			contestTeam.POST("/leave", middleware.ContestIsComing, LeaveTeam)
		}

		// 比赛公告
		{
			contest.GET("/notices", GetNotices)
		}

		contest.GET("/challenges", middleware.CheckVerified, middleware.SetTeamByUser, middleware.CheckBanned, middleware.ContestIsNotComing, GetContestChallenges)
		contest.GET("/challenges/categories", middleware.CheckVerified, middleware.SetTeamByUser, middleware.CheckBanned, middleware.ContestIsNotComing, GetContestChallengeCategories)
		contestChallenge := contest.Group(
			"/challenges/:challengeID",
			middleware.CheckVerified, middleware.SetTeamByUser, middleware.CheckBanned, middleware.ContestIsNotComing, middleware.SetContestChallenge,
		)
		{
			contestChallenge.GET("", GetContestChallengeStatus)
			contestChallenge.POST("/init",
				middleware.RateLimit("init_flag", 1, time.Minute),
				middleware.ContestIsRunning, middleware.CheckSolved, InitTeamFlag,
			)
			contestChallenge.GET("/attachment",
				middleware.RateLimit("download_attachment", 10, time.Minute),
				middleware.SetAttachmentFile(false), DownloadFile(model.DownloadAttachmentEventType),
			)
			contestChallenge.POST("/reset",
				middleware.RateLimit("init_flag", 1, time.Minute),
				middleware.ContestIsRunning, middleware.CheckIfGenerated, middleware.CheckSolved, ResetTeamFlag,
			)
			contestChallenge.POST("/start",
				middleware.CheckChallengeType(model.PodsChallengeType),
				middleware.RateLimit("start_victim", 1, time.Minute), middleware.CheckTeamVictimCount, middleware.CheckIfGenerated, StartVictim,
			)
			contestChallenge.POST("/increase",
				middleware.CheckChallengeType(model.PodsChallengeType), middleware.ContestIsRunning,
				middleware.CheckIfGenerated, IncreaseVictimDuration,
			)
			contestChallenge.POST("/stop", middleware.CheckChallengeType(model.PodsChallengeType), middleware.CheckIfGenerated, StopVictim)
			contestChallenge.POST("/submit",
				middleware.RateLimit("submit_flag", 1, time.Second),
				middleware.ContestIsRunning, middleware.CheckIfGenerated, middleware.CheckSolved, SubmitFlag,
			)
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

	admin := auth.Group("/admin", middleware.CheckRole(true))
	{
		admin.GET("/me", GetAdmin)
		admin.PUT("/me/password", AdminChangePassword)
		admin.PUT("/me", UpdateAdmin)
		admin.POST("/me/avatar", UploadAvatar("admin"))
		admin.POST("", CreateAdmin)

		admin.GET("/search", Search)

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

		admin.GET("/oauth", GetOauthProviders)
		admin.POST("/oauth", CreateOauthProvider)
		adminOauth := admin.Group("/oauth/:oauthID", middleware.SetOauth)
		{
			adminOauth.GET("", GetOauthProvider)
			adminOauth.PUT("", UpdateOauthProvider)
			adminOauth.POST("/avatar", UploadAvatar("oauth"))
			adminOauth.DELETE("", DeleteOauthProvider)
		}

		admin.GET("/email", GetEmails)
		admin.GET("/smtp", GetSmtps)
		admin.POST("/smtp", CreateSmtp)
		adminSmtp := admin.Group("/smtp/:smtpID", middleware.SetSmtp)
		{
			adminSmtp.GET("", GetSmtp)
			adminSmtp.PUT("", UpdateSmtp)
			adminSmtp.DELETE("", DeleteSmtp)

			adminSmtp.GET("/email", GetEmails)
		}

		admin.GET("/webhook", GetWebhooks)
		admin.GET("/webhook/events", GetEventTypes)
		admin.GET("/webhook/history", GetWebhookHistory)
		admin.POST("/webhook", CreateWebhook)
		adminWebhook := admin.Group("/webhook/:webhookID", middleware.SetWebhook)
		{
			adminWebhook.GET("", GetWebhook)
			adminWebhook.PUT("", UpdateWebhook)
			adminWebhook.DELETE("", DeleteWebhook)

			adminWebhook.GET("/history", GetWebhookHistory)
		}

		admin.GET("/challenges", GetChallenges)
		admin.GET("/challenges/categories", GetChallengeCategories)
		admin.POST("/challenges", CreateChallenge)
		adminChallenge := admin.Group("/challenges/:challengeID", middleware.SetChallenge)
		{
			adminChallenge.GET("", GetChallenge)
			adminChallenge.GET("/download", middleware.SetChallengeFile, DownloadFile(model.DownloadAttachmentEventType))
			adminChallenge.PUT("", UpdateChallenge)
			adminChallenge.DELETE("", DeleteChallenge)
			adminChallenge.POST("/upload", UploadChallengeFile)

			adminChallengeTest := adminChallenge.Group("/test")
			{
				adminChallengeTest.GET("", GetTestChallengeStatus)
				adminChallengeTest.GET("/attachment",
					middleware.RateLimit("download_attachment", 10, time.Minute),
					middleware.SetAttachmentFile(true), DownloadFile(model.DownloadAttachmentEventType),
				)
				adminChallengeTest.POST("/start",
					middleware.CheckChallengeType(model.PodsChallengeType),
					middleware.RateLimit("start_victim", 10, time.Minute), StartTestVictim,
				)
				adminChallengeTest.POST("/stop", middleware.CheckChallengeType(model.PodsChallengeType), StopTestVictim)
			}
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
			adminContest.GET("/scoreboard", GetScoreboard)
			adminContest.GET("/timeline", GetRankTimeline)

			adminContest.GET("/teams", GetTeams)
			adminContestTeam := adminContest.Group("/teams/:teamID", middleware.SetTeam)
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
					adminTraffic := adminContainer.Group("/traffic")
					{
						adminTraffic.GET("/download", middleware.SetTrafficFile, DownloadFile(model.DownloadTrafficEventType))
						adminTraffic.GET("", GetTraffics)
					}
				}

				adminContestTeam.GET("/writeups", GetWriteUPs)
				adminContestTeam.GET("/writeups/:fileID", middleware.SetFile(model.WriteUPFileType), DownloadFile(model.DownloadWriteUpEventType))
			}

			adminContest.GET("/notices", GetNotices)
			adminContest.POST("/notices", CreateNotice)
			adminContestNotice := adminContest.Group("/notices/:noticeID", middleware.SetNotice)
			{
				adminContestNotice.PUT("", UpdateNotice)
				adminContestNotice.DELETE("", DeleteNotice)
			}

			adminContest.GET("/cheats", GetCheats)
			adminContestCheat := adminContest.Group("/cheats/:cheatID", middleware.SetCheat)
			{
				adminContestCheat.PUT("", UpdateCheat)
				//adminContestCheat.DELETE("")
			}

			adminContest.GET("/challenges", GetContestChallenges)
			adminContest.GET("/challenges/categories", GetContestChallengeCategories)
			adminContest.POST("/challenges", AddContestChallenge)
			adminContestChallenge := adminContest.Group("/challenges/:challengeID", middleware.SetContestChallenge)
			{
				adminContestChallenge.GET("", GetContestChallenge)
				adminContestChallenge.PUT("", UpdateContestChallenge)
				adminContestChallenge.DELETE("", DeleteContestChallenge)

				//不允许后期创建和删除
				adminContestChallenge.GET("/flags", GetContestFlags)
				adminContestFlag := adminContestChallenge.Group("/flags/:flagID", middleware.SetContestFlag)
				{
					adminContestFlag.GET("", GetContestFlag)
					adminContestFlag.PUT("", UpdateContestFlag)
				}
			}

			adminContestImages := adminContest.Group("/images")
			{
				adminContestImages.GET("", GetContestChallengeImage)
				adminContestImages.POST("", WarmUpContestChallengeImage)
			}

			adminContestVictim := adminContest.Group("/victims")
			{
				adminContestVictim.GET("", GetContestVictims)
				adminContestVictim.POST("", StartContestVictims)
				adminContestVictim.DELETE("", StopContestVictims)
			}
		}

		admin.GET("/files", GetFiles)
		admin.DELETE("/files", DeleteFiles)
		admin.GET("/files/:fileID", middleware.SetFile(""), DownloadFile(model.DownloadFileEventType))

		admin.GET("/logs", GetLogs)
	}
	return router
}
