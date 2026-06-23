package router

import (
	"CBCTF/frontend"
	"CBCTF/internal/config"
	"CBCTF/internal/log"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

func Init() *gin.Engine {
	pictureMaxBytes := int64(config.Env.Gin.Upload.Picture) << 20
	challengeMaxBytes := int64(config.Env.Gin.Upload.Challenge) << 20
	writeupMaxBytes := int64(config.Env.Gin.Upload.Writeup) << 20

	gin.SetMode(strings.ToLower(config.Env.Gin.Mode))
	router := gin.New()

	log.Logger.Infof("Trust proxies: %s", config.Env.Gin.Proxies)
	if err := router.SetTrustedProxies(config.Env.Gin.Proxies); err != nil {
		log.Logger.Warningf("Set trusted proxies failed: %s", err)
	}

	router.Use(middleware.SetTrace, middleware.SetMagic, middleware.Cors())

	{
		router.GET("/", func(ctx *gin.Context) {
			ctx.Redirect(http.StatusFound, fmt.Sprintf("%s/platform", config.Env.Host))
		})
		router.Use(func(ctx *gin.Context) {
			if strings.HasPrefix(ctx.Request.URL.Path, "/platform") {
				ctx.Header("Cache-Control", "public, max-age=31536000, immutable")
			}
			ctx.Next()
		})
		router.StaticFS("/platform", http.FS(frontend.SubFS))
	}

	router.Use(
		middleware.Logger, middleware.Prometheus, middleware.AccessLog, middleware.Events, middleware.Recovery,
		middleware.RateLimit("globals", config.Env.Gin.RateLimit.Global, time.Minute),
	)

	RegisterMetricsRouter(router)

	{
		router.POST("/register", middleware.RateLimit("register", 1, time.Minute), Register)
		router.POST("/login", middleware.RateLimit("login", 10, time.Minute), Login)
		router.DELETE("/logout", Logout)

		router.POST("/password/forgot", middleware.RateLimit("forgot-password", 1, time.Minute), ForgotPassword)
		router.POST("/password/reset", middleware.RateLimit("reset-password", 1, time.Minute), ResetPassword)

		router.GET("/config", PublicSystemConfig)

		RegisterOauthRouter()
		router.GET("/oauth", ListOauth)
		router.GET("/oauth/token", ExchangeOauthCode)
		oauth := router.Group("/oauth/:oauth", middleware.SetOauthUri)
		{
			oauth.GET("", Oauth)
			oauth.GET("/callback", OauthCallback)
		}

		router.POST("/verify", middleware.RateLimit("verify", 5, time.Minute), VerifyEmail)
		router.GET("/assets", DefaultAssets)
		router.GET("/pictures/:fileID", middleware.SetFile(model.PictureFileType), DownloadFile(model.SkipEventType))

		router.GET("/branding", GetBranding)
		router.GET("/stats", HomePage)
		router.GET("/contests", GetContests)
	}

	if strings.ToLower(config.Env.Gin.Mode) != gin.ReleaseMode {
		pprof.Register(router)
	}

	auth := router.Group("", middleware.CheckAuth, middleware.CheckPermission)

	user := auth.Group("/me")
	{
		user.GET("", GetUser)
		user.GET("/permissions", GetAccessibleRoutes)
		user.PUT("/password", ChangePwd)
		user.PUT("", UpdateUser)
		user.DELETE("", DeleteUser)
		user.POST("/picture", middleware.LimitUploadSize(pictureMaxBytes), UploadPicture("self"))
		user.POST("/activate",
			middleware.RateLimit("activate", 1, time.Minute), ActivateEmail,
		)
	}

	contest := auth.Group("/contests/:contestID", middleware.SetContest)
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
			contestTeam.PUT("/captcha", middleware.CheckCaptain, UpdateCaptcha)
			contestTeam.PUT("", middleware.CheckCaptain, UpdateTeam)
			contestTeam.POST("/picture", middleware.CheckCaptain, middleware.LimitUploadSize(pictureMaxBytes), UploadPicture("team"))
			contestTeam.DELETE("", middleware.ContestIsComing, middleware.CheckCaptain, DeleteTeam)
			contestTeam.POST("/kick", middleware.ContestIsComing, middleware.CheckCaptain, KickMember)
			contestTeam.POST("/leave", middleware.ContestIsComing, LeaveTeam)
		}

		// 比赛公告
		{
			contest.GET("/notices", GetNotices)
		}

		contest.GET("/challenges",
			middleware.CheckVerified, middleware.SetTeamByUser, middleware.CheckBanned, middleware.ContestIsNotComing, GetContestChallenges,
		)
		contest.GET("/challenges/categories",
			middleware.CheckVerified, middleware.SetTeamByUser, middleware.CheckBanned, middleware.ContestIsNotComing, GetContestChallengeCategories,
		)
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
			contestChallenge.POST("/extend",
				middleware.CheckChallengeType(model.PodsChallengeType), middleware.ContestIsRunning,
				middleware.CheckIfGenerated, ExtendVictimDuration,
			)
			contestChallenge.POST("/stop",
				middleware.CheckChallengeType(model.PodsChallengeType), middleware.CheckIfGenerated, StopVictim,
			)
			contestChallenge.POST("/submit",
				middleware.RateLimit("submit_flag", 1, time.Second),
				middleware.ContestIsRunning, middleware.CheckIfGenerated, middleware.CheckSolved, SubmitFlag,
			)
		}

		// Writeup
		contestWriteup := contest.Group(
			"/writeups",
			middleware.CheckVerified, middleware.SetTeamByUser, middleware.CheckBanned, middleware.ContestIsNotComing,
		)
		{
			contestWriteup.POST("", middleware.LimitUploadSize(writeupMaxBytes), UploadWriteup)
			contestWriteup.GET("", GetWriteUPs)
		}
	}

	admin := auth.Group("/admin", middleware.SetFullAccess)
	{
		admin.GET("/ip", SearchIP)
		admin.GET("/models", GetAllowQueryModels)
		admin.GET("/search", Search)

		// 系统管理
		adminSystem := admin.Group("/system")
		{
			adminSystem.GET("/status", SystemStatus)
			adminSystem.GET("/config", SystemConfig)
			adminSystem.PUT("/config", UpdateSystem)
			adminSystem.POST("/geocity-db", UploadGeoCityDB)
			adminSystem.POST("/restart", RestartSystem)
		}

		admin.GET("/branding", GetAdminBranding)
		admin.PUT("/branding", UpdateBranding)
		admin.POST("/branding/logo", middleware.LimitUploadSize(pictureMaxBytes), UploadPicture("branding"))

		admin.GET("/permissions", GetPermissions)
		adminPermission := admin.Group("/permissions/:permissionID", middleware.SetPermission)
		{
			adminPermission.PUT("", UpdatePermission)
		}

		admin.GET("/roles", GetRoles)
		admin.POST("/roles", CreateRole)
		adminRole := admin.Group("/roles/:roleID", middleware.SetRole)
		{
			adminRole.GET("", GetRole)
			adminRole.GET("/permissions", GetRolePermissions)
			adminRole.PUT("", UpdateRole)
			adminRole.DELETE("", DeleteRole)
			adminRole.POST("/permissions", AssignPermission)
			adminRole.DELETE("/permissions", RevokePermission)
		}

		admin.GET("/groups", GetGroups)
		admin.POST("/groups", CreateGroup)
		adminGroup := admin.Group("/groups/:groupID", middleware.SetGroup)
		{
			adminGroup.GET("", GetGroup)
			adminGroup.GET("/users", GetGroupUsers)
			adminGroup.GET("/users/available", GetGroupAvailableUsers)
			adminGroup.PUT("", UpdateGroup)
			adminGroup.DELETE("", DeleteGroup)
			adminGroup.POST("/users", AssignUserToGroup)
			adminGroup.DELETE("/users", RemoveUserFromGroup)
		}

		admin.GET("/users", GetUsers)
		admin.POST("/users", CreateUser)
		adminUser := admin.Group("/users/:userID", middleware.SetUser)
		{
			adminUser.GET("", GetUser)
			adminUser.PUT("", UpdateUser)
			adminUser.DELETE("", DeleteUser)
			adminUser.POST("/picture", middleware.LimitUploadSize(pictureMaxBytes), UploadPicture("user"))
		}

		admin.GET("/oauth", GetOauthProviders)
		admin.POST("/oauth", CreateOauthProvider)
		adminOauth := admin.Group("/oauth/:oauthID", middleware.SetOauth)
		{
			adminOauth.GET("", GetOauthProvider)
			adminOauth.PUT("", UpdateOauthProvider)
			adminOauth.POST("/picture", middleware.LimitUploadSize(pictureMaxBytes), UploadPicture("oauth"))
			adminOauth.DELETE("", DeleteOauthProvider)
		}

		admin.GET("/email", GetEmails)
		admin.GET("/smtp", GetSmtps)
		admin.POST("/smtp", CreateSmtp)
		adminSmtp := admin.Group("/smtp/:smtpID", middleware.SetSmtp)
		{
			adminSmtp.GET("", ReadSmtp)
			adminSmtp.PUT("", UpdateSmtp)
			adminSmtp.DELETE("", DeleteSmtp)
			adminSmtp.POST("/test", TestSmtp)

			adminSmtp.GET("/email", GetEmails)
		}

		admin.GET("/cronjobs", GetCronJobs)
		adminCronJob := admin.Group("/cronjobs/:cronJobID", middleware.SetCronJob)
		{
			adminCronJob.GET("", GetCronJob)
			adminCronJob.PUT("", UpdateCronJob)
		}

		admin.GET("/webhook", GetWebhooks)
		admin.GET("/webhook/events", GetEventTypes)
		admin.GET("/webhook/history", GetWebhookHistory)
		admin.POST("/webhook", CreateWebhook)
		adminWebhook := admin.Group("/webhook/:webhookID", middleware.SetWebhook)
		{
			adminWebhook.GET("", ReadWebhook)
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
			adminChallenge.GET("/download",
				middleware.SetChallengeFile, DownloadFile(model.DownloadAttachmentEventType),
			)
			adminChallenge.PUT("", UpdateChallenge)
			adminChallenge.DELETE("", DeleteChallenge)
			adminChallenge.POST("/upload", middleware.LimitUploadSize(challengeMaxBytes), UploadChallenge)

			adminChallengeTest := adminChallenge.Group("/test")
			{
				adminChallengeTest.GET("", GetTestChallengeStatus)
				adminChallengeTest.GET("/attachment",
					middleware.RateLimit("download_attachment", 10, time.Minute),
					middleware.SetAttachmentFile(true), DownloadFile(model.DownloadAttachmentEventType),
				)
				adminChallengeTest.POST("/start",
					middleware.CheckChallengeType(model.PodsChallengeType),
					middleware.RateLimit("test_victim", 10, time.Minute), StartTestVictim,
				)
				adminChallengeTest.POST("/stop", middleware.CheckChallengeType(model.PodsChallengeType), StopTestVictim)
			}
		}

		adminVictim := admin.Group("/victims")
		{
			adminVictim.GET("", GetVictims)
			adminVictim.DELETE("", StopVictims)
			adminSingleVictim := adminVictim.Group("/:victimID", middleware.SetVictim)
			{
				adminTraffic := adminSingleVictim.Group("/traffic")
				{
					adminTraffic.GET("/download", middleware.SetTrafficFile, DownloadFile(model.DownloadTrafficEventType))
					adminTraffic.GET("", GetTraffics)
				}
				adminSingleVictim.GET("/pods", GetVictimPods)
				adminSingleVictim.GET("/pods/logs", GetVictimPodLogs)
			}
		}

		adminGenerator := admin.Group("/generators")
		{
			adminGenerator.GET("", GetGenerators)
			adminGenerator.POST("", middleware.RateLimit("test_generators", 1, time.Minute), StartGenerator)
			adminGenerator.DELETE("", StopGenerator)
			adminGenerator.GET("/:generatorID/logs", middleware.SetGenerator, GetGeneratorLogs)
		}

		adminImages := admin.Group("/images")
		{
			adminImages.GET("", GetImages)
			adminImages.POST("", middleware.RateLimit("warmup_images", 1, time.Minute), PullImages)
		}

		admin.GET("/contests", GetContests)
		admin.POST("/contests", CreateContest)
		adminContest := admin.Group("/contests/:contestID", middleware.SetContest)
		{
			adminContest.GET("", GetContest)
			adminContest.PUT("", UpdateContest)
			adminContest.DELETE("", DeleteContest)
			adminContest.POST("/picture", middleware.LimitUploadSize(pictureMaxBytes), UploadPicture("contest"))
			adminContest.GET("/rank", GetTeamRanking)
			adminContest.GET("/scoreboard", GetScoreboard)
			adminContest.GET("/timeline", GetRankTimeline)
			adminContest.GET("/writeups/export", ExportContestWriteups)

			adminContest.GET("/teams", GetTeams)
			adminContestTeam := adminContest.Group("/teams/:teamID", middleware.SetTeam)
			{
				adminContestTeam.GET("", GetTeam)
				adminContestTeam.GET("/users", GetTeammates)
				adminContestTeam.PUT("", UpdateTeam)
				adminContestTeam.DELETE("", DeleteTeam)
				adminContestTeam.POST("/kick", KickMember)
				adminContestTeam.POST("/picture", middleware.LimitUploadSize(pictureMaxBytes), UploadPicture("team"))

				adminContestTeam.GET("/flags", GetTeamFlags)

				adminContestTeam.GET("/submissions", GetSubmissions)

				adminContestTeam.GET("/victims", GetVictimHistories)
				adminContestVictim := adminContestTeam.Group("/victims/:victimID", middleware.SetVictim)
				{
					adminContestVictimTraffic := adminContestVictim.Group("/traffic")
					{
						adminContestVictimTraffic.GET("/download", middleware.SetTrafficFile, DownloadFile(model.DownloadTrafficEventType))
						adminContestVictimTraffic.GET("", GetTraffics)
					}
				}

				adminContestTeam.GET("/writeups", GetWriteUPs)
				adminContestTeam.GET("/writeups/:fileID",
					middleware.SetTeamWriteupFile, DownloadFile(model.DownloadWriteupEventType),
				)
			}

			adminContest.GET("/notices", GetNotices)
			adminContest.POST("/notices", CreateNotice)
			adminContestNotice := adminContest.Group("/notices/:noticeID", middleware.SetNotice)
			{
				adminContestNotice.PUT("", UpdateNotice)
				adminContestNotice.DELETE("", DeleteNotice)
			}

			adminContest.GET("/cheats", GetCheats)
			adminContest.DELETE("/cheats", DeleteCheat(true))
			adminContest.POST("/cheats", CheckCheat)
			adminContestCheat := adminContest.Group("/cheats/:cheatID", middleware.SetCheat)
			{
				adminContestCheat.PUT("", UpdateCheat)
				adminContestCheat.DELETE("", DeleteCheat(false))
			}

			adminContest.GET("/challenges", GetAllContestChallenges)
			adminContest.GET("/challenges/others", GetChallengeNotInContest)
			adminContest.GET("/challenges/categories", GetContestChallengeCategories)
			adminContest.POST("/challenges", AddContestChallenge)
			adminContestChallenge := adminContest.Group("/challenges/:challengeID", middleware.SetContestChallenge)
			{
				adminContestChallenge.GET("", GetAdminContestChallenge)
				adminContestChallenge.PUT("", UpdateContestChallenge)
				adminContestChallenge.DELETE("", DeleteContestChallenge)

				//不允许后期创建和删除
				adminContestChallenge.GET("/flags", GetContestFlags)
				adminContestFlag := adminContestChallenge.Group("/flags/:flagID", middleware.SetContestFlag)
				{
					adminContestFlag.GET("", ReadContestFlag)
					adminContestFlag.PUT("", UpdateContestFlag)
					adminContestFlag.GET("/solvers", GetContestFlagSolvers)
				}
			}

			adminContestImages := adminContest.Group("/images")
			{
				adminContestImages.GET("", GetContestChallengeImage)
				adminContestImages.POST("", middleware.RateLimit("warmup_images", 1, time.Minute), PullImages)
			}

			adminContestVictim := adminContest.Group("/victims")
			{
				adminContestVictim.GET("", GetVictims)
				adminContestVictim.POST("", middleware.RateLimit("warmup_victims", 1, time.Minute), StartVictims)
				adminContestVictim.DELETE("", StopVictims)
				adminContestSingleVictim := adminContestVictim.Group("/:victimID", middleware.SetVictim)
				{
					adminContestSingleVictim.GET("/pods", GetVictimPods)
					adminContestSingleVictim.GET("/pods/logs", GetVictimPodLogs)
				}
			}

			adminContestGenerator := adminContest.Group("/generators")
			{
				adminContestGenerator.GET("", GetGenerators)
				adminContestGenerator.POST("", middleware.RateLimit("warmup_generators", 1, time.Minute), StartGenerator)
				adminContestGenerator.DELETE("", StopGenerator)
				adminContestGenerator.GET("/:generatorID/logs", middleware.SetGenerator, GetGeneratorLogs)
			}
		}

		admin.GET("/files", GetFiles)
		admin.DELETE("/files", DeleteFiles)
		admin.GET("/files/:fileID", middleware.SetFile(""), DownloadFile(model.DownloadFileEventType))

		admin.GET("/tasks", GetTasks)
		admin.GET("/tasks/live", GetLiveTasks)

		admin.GET("/logs", GetLogs)
	}
	return router
}
