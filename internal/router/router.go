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
	gin.SetMode(strings.ToLower(config.Env.Gin.Mode))
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
			ctx.Redirect(http.StatusFound, fmt.Sprintf("%s/platform", config.Env.Host))
		})
		router.StaticFS("/platform", http.FS(frontend.SubFS))
	}

	router.Use(
		middleware.SetTrace, middleware.SetMagic, middleware.Logger, middleware.Prometheus, middleware.AccessLog,
		middleware.I18n, middleware.RateLimit("globals", config.Env.Gin.RateLimit.Global, time.Minute), middleware.Events,
	)

	if strings.ToLower(config.Env.Gin.Mode) != gin.ReleaseMode {
		pprof.Register(router)
	}

	RegisterMetricsRouter(router)

	{
		router.POST("/register", middleware.RateLimit("register", 1, time.Minute), Register)
		router.POST("/login", Login)

		RegisterOauthRouter()
		router.GET("/oauth", ListOauth)
		router.GET("/oauth/token", ExchangeOauthCode)
		oauth := router.Group("/oauth/:oauth", middleware.SetOauthUri)
		{
			oauth.GET("", Oauth)
			oauth.GET("/callback", OauthCallback)
		}

		router.GET("/verify", VerifyEmail)
		router.GET("/assets", DefaultAssets)
		router.GET("/pictures/:fileID", middleware.SetFile(model.PictureFileType), DownloadFile(model.SkipEventType))

		router.GET("/stats", HomePage)
		router.GET("/contests", GetContests)
	}

	auth := router.Group("", middleware.CheckAuth, middleware.CheckPermission)

	user := auth.Group("/me")
	{
		user.GET("", GetUser)
		user.GET("/permissions", GetAccessibleRoutes)
		user.PUT("/password", ChangePwd)
		user.PUT("", UpdateUser)
		user.DELETE("", DeleteUser)
		user.POST("/picture", UploadPicture("self"))
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
		contest.POST("/teams/join",
			middleware.ContestIsNotOver, middleware.CheckVerified, JoinTeam,
		)
		contest.POST("/teams/create",
			middleware.ContestIsNotOver, middleware.CheckVerified, CreateTeam,
		)

		contestTeam := contest.Group("/teams/me", middleware.CheckVerified, middleware.SetTeamByUser)
		{
			contestTeam.GET("", GetTeam)
			contestTeam.GET("/captcha", GetTeamCaptcha)
			contestTeam.GET("/users", GetTeammates)
			contestTeam.PUT("/captcha",
				middleware.ContestIsNotOver, middleware.CheckCaptain, UpdateCaptcha,
			)
			contestTeam.PUT("",
				middleware.ContestIsNotOver, middleware.CheckCaptain, UpdateTeam,
			)
			contestTeam.POST("/picture",
				middleware.ContestIsNotOver, middleware.CheckCaptain, UploadPicture("team"),
			)
			contestTeam.DELETE("",
				middleware.ContestIsComing, middleware.CheckCaptain, DeleteTeam,
			)
			contestTeam.POST("/kick",
				middleware.ContestIsComing, middleware.CheckCaptain, KickMember,
			)
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
			contestChallenge.POST("/increase",
				middleware.CheckChallengeType(model.PodsChallengeType), middleware.ContestIsRunning,
				middleware.CheckIfGenerated, IncreaseVictimDuration,
			)
			contestChallenge.POST("/stop",
				middleware.CheckChallengeType(model.PodsChallengeType), middleware.CheckIfGenerated, StopVictim,
			)
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
			adminSystem.POST("/restart", RestartSystem)
		}

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
			adminUser.POST("/picture", UploadPicture("user"))
		}

		admin.GET("/oauth", GetOauthProviders)
		admin.POST("/oauth", CreateOauthProvider)
		adminOauth := admin.Group("/oauth/:oauthID", middleware.SetOauth)
		{
			adminOauth.PUT("", UpdateOauthProvider)
			adminOauth.POST("/picture", UploadPicture("oauth"))
			adminOauth.DELETE("", DeleteOauthProvider)
		}

		admin.GET("/email", GetEmails)
		admin.GET("/smtp", GetSmtps)
		admin.POST("/smtp", CreateSmtp)
		adminSmtp := admin.Group("/smtp/:smtpID", middleware.SetSmtp)
		{
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
			adminWebhook.PUT("", UpdateWebhook)
			adminWebhook.DELETE("", DeleteWebhook)

			adminWebhook.GET("/history", GetWebhookHistory)
		}

		admin.GET("/challenges", GetChallenges)
		admin.GET("/challenges/categories", GetChallengeCategories)
		admin.POST("/challenges", CreateChallenge)
		adminChallenge := admin.Group("/challenges/:challengeID", middleware.SetChallenge)
		{
			adminChallenge.GET("/download",
				middleware.SetChallengeFile, DownloadFile(model.DownloadAttachmentEventType),
			)
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
			adminContest.POST("/picture", UploadPicture("contest"))
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
				adminContestTeam.POST("/picture", UploadPicture("team"))

				adminContestTeam.GET("/flags", GetTeamFlags)

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

				adminContestTeam.GET("/writeups",
					GetWriteUPs,
				)
				adminContestTeam.GET("/writeups/:fileID",
					middleware.SetFile(model.WriteupFileType), DownloadFile(model.DownloadWriteUpEventType),
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

			adminContest.GET("/challenges", GetContestChallenges)
			adminContest.GET("/challenges/others", GetChallengeNotInContest)
			adminContest.GET("/challenges/categories", GetContestChallengeCategories)
			adminContest.POST("/challenges", AddContestChallenge)
			adminContestChallenge := adminContest.Group("/challenges/:challengeID", middleware.SetContestChallenge)
			{
				adminContestChallenge.PUT("", UpdateContestChallenge)
				adminContestChallenge.DELETE("", DeleteContestChallenge)

				//不允许后期创建和删除
				adminContestChallenge.GET("/flags", GetContestFlags)
				adminContestFlag := adminContestChallenge.Group("/flags/:flagID", middleware.SetContestFlag)
				{
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
