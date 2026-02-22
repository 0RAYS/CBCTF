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
		router.GET("/pictures/:fileID", middleware.SetFile(model.PictureFileType), DownloadFile(model.SkipEventType))

		router.GET("/stats", HomePage)
		router.GET("/contests", GetContests)
	}

	auth := router.Group("", middleware.CheckAuth)

	user := auth.Group("/me", middleware.CheckRole(false))
	{
		user.GET("", middleware.RequirePermission(model.PermSelfRead), GetUser)
		user.PUT("/password", middleware.RequirePermission(model.PermSelfUpdate), ChangePwd)
		user.PUT("", middleware.RequirePermission(model.PermSelfUpdate), UpdateUser)
		user.DELETE("", middleware.RequirePermission(model.PermSelfDelete), DeleteUser)
		user.POST("/picture", middleware.RequirePermission(model.PermSelfUpdate), UploadPicture("self-user"))
		user.POST("/activate",
			middleware.RequirePermission(model.PermSelfActivate),
			middleware.RateLimit("activate", 1, time.Minute), ActivateEmail,
		)
	}

	contest := auth.Group("/contests/:contestID", middleware.CheckRole(false), middleware.SetContest)
	{
		contest.GET("", middleware.RequirePermission(model.PermUserContestRead), GetContest)
		contest.GET("/rank", middleware.RequirePermission(model.PermUserContestRank), GetTeamRanking)
		contest.GET("/scoreboard", middleware.RequirePermission(model.PermUserContestRank), GetScoreboard)
		contest.GET("/timeline", middleware.RequirePermission(model.PermUserContestRank), GetRankTimeline)
		contest.POST("/teams/join", middleware.RequirePermission(model.PermUserTeamJoin),
			middleware.ContestIsNotOver, middleware.CheckVerified, JoinTeam,
		)
		contest.POST("/teams/create", middleware.RequirePermission(model.PermUserTeamCreate),
			middleware.ContestIsNotOver, middleware.CheckVerified, CreateTeam,
		)

		contestTeam := contest.Group("/teams/me", middleware.CheckVerified, middleware.SetTeamByUser)
		{
			contestTeam.GET("", middleware.RequirePermission(model.PermUserTeamRead), GetTeam)
			contestTeam.GET("/captcha", middleware.RequirePermission(model.PermUserTeamRead), GetTeamCaptcha)
			contestTeam.GET("/users", middleware.RequirePermission(model.PermUserTeamRead), GetTeammates)
			contestTeam.PUT("/captcha", middleware.RequirePermission(model.PermUserTeamUpdate),
				middleware.ContestIsNotOver, middleware.CheckCaptain, UpdateCaptcha,
			)
			contestTeam.PUT("", middleware.RequirePermission(model.PermUserTeamUpdate),
				middleware.ContestIsNotOver, middleware.CheckCaptain, UpdateTeam,
			)
			contestTeam.POST("/picture", middleware.RequirePermission(model.PermUserTeamUpdate),
				middleware.ContestIsNotOver, middleware.CheckCaptain, UploadPicture("team"),
			)
			contestTeam.DELETE("", middleware.RequirePermission(model.PermUserTeamDelete),
				middleware.ContestIsComing, middleware.CheckCaptain, DeleteTeam,
			)
			contestTeam.POST("/kick", middleware.RequirePermission(model.PermUserTeamUpdate),
				middleware.ContestIsComing, middleware.CheckCaptain, KickMember,
			)
			contestTeam.POST("/leave", middleware.RequirePermission(model.PermUserTeamRead), middleware.ContestIsComing, LeaveTeam)
		}

		// 比赛公告
		{
			contest.GET("/notices", middleware.RequirePermission(model.PermUserNoticeList), GetNotices)
		}

		contest.GET("/challenges", middleware.RequirePermission(model.PermUserChallengeList),
			middleware.CheckVerified, middleware.SetTeamByUser, middleware.CheckBanned, middleware.ContestIsNotComing, GetContestChallenges,
		)
		contest.GET("/challenges/categories", middleware.RequirePermission(model.PermUserChallengeList),
			middleware.CheckVerified, middleware.SetTeamByUser, middleware.CheckBanned, middleware.ContestIsNotComing, GetContestChallengeCategories,
		)
		contestChallenge := contest.Group(
			"/challenges/:challengeID",
			middleware.CheckVerified, middleware.SetTeamByUser, middleware.CheckBanned, middleware.ContestIsNotComing, middleware.SetContestChallenge,
		)
		{
			contestChallenge.GET("", middleware.RequirePermission(model.PermUserChallengeRead), GetContestChallengeStatus)
			contestChallenge.POST("/init",
				middleware.RequirePermission(model.PermUserChallengeInit),
				middleware.RateLimit("init_flag", 1, time.Minute),
				middleware.ContestIsRunning, middleware.CheckSolved, InitTeamFlag,
			)
			contestChallenge.GET("/attachment",
				middleware.RequirePermission(model.PermUserChallengeRead),
				middleware.RateLimit("download_attachment", 10, time.Minute),
				middleware.SetAttachmentFile(false), DownloadFile(model.DownloadAttachmentEventType),
			)
			contestChallenge.POST("/reset",
				middleware.RequirePermission(model.PermUserChallengeReset),
				middleware.RateLimit("init_flag", 1, time.Minute),
				middleware.ContestIsRunning, middleware.CheckIfGenerated, middleware.CheckSolved, ResetTeamFlag,
			)
			contestChallenge.POST("/start",
				middleware.RequirePermission(model.PermUserVictimControl),
				middleware.CheckChallengeType(model.PodsChallengeType),
				middleware.RateLimit("start_victim", 1, time.Minute), middleware.CheckTeamVictimCount, middleware.CheckIfGenerated, StartVictim,
			)
			contestChallenge.POST("/increase",
				middleware.RequirePermission(model.PermUserVictimControl),
				middleware.CheckChallengeType(model.PodsChallengeType), middleware.ContestIsRunning,
				middleware.CheckIfGenerated, IncreaseVictimDuration,
			)
			contestChallenge.POST("/stop", middleware.RequirePermission(model.PermUserVictimControl),
				middleware.CheckChallengeType(model.PodsChallengeType), middleware.CheckIfGenerated, StopVictim,
			)
			contestChallenge.POST("/submit",
				middleware.RequirePermission(model.PermUserChallengeSubmit),
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
			contestWriteUp.POST("", middleware.RequirePermission(model.PermUserWriteupUpload), UploadWriteUp)
			contestWriteUp.GET("", middleware.RequirePermission(model.PermUserWriteupList), GetWriteUPs)
		}
	}

	admin := auth.Group("/admin", middleware.CheckRole(true))
	{
		admin.GET("/ip", middleware.RequirePermission(model.PermAdminIPSearch), SearchIP)

		admin.GET("/me", GetAdmin)
		admin.PUT("/me/password", AdminChangePassword)
		admin.PUT("/me", UpdateAdmin)
		admin.POST("/me/picture", UploadPicture("admin"))
		admin.POST("", CreateAdmin)

		admin.GET("/models", middleware.RequirePermission(model.PermAdminModelsSearch), GetAllowQueryModels)
		admin.GET("/search", middleware.RequirePermission(model.PermAdminModelsSearch), Search)

		// 系统管理
		adminSystem := admin.Group("/system")
		{
			adminSystem.GET("/status", middleware.RequirePermission(model.PermAdminSystemRead), SystemStatus)
			adminSystem.GET("/config", middleware.RequirePermission(model.PermAdminSystemRead), SystemConfig)
			adminSystem.PUT("/config", middleware.RequirePermission(model.PermAdminSystemUpdate), UpdateSystem)
			adminSystem.POST("/restart", middleware.RequirePermission(model.PermAdminSystemRestart), RestartSystem)
		}

		admin.GET("/permissions", middleware.RequirePermission(model.PermAdminPermissionList), GetPermissions)
		adminPermission := admin.Group("/permissions/:permissionID", middleware.SetPermission)
		{
			adminPermission.GET("", middleware.RequirePermission(model.PermAdminPermissionRead), GetPermission)
			adminPermission.PUT("", middleware.RequirePermission(model.PermAdminPermissionUpdate), UpdatePermission)
		}

		admin.GET("/roles", middleware.RequirePermission(model.PermAdminRoleList), GetRoles)
		admin.POST("/roles", middleware.RequirePermission(model.PermAdminRoleCreate), CreateRole)
		adminRole := admin.Group("/roles/:roleID", middleware.SetRole)
		{
			adminRole.GET("", middleware.RequirePermission(model.PermAdminRoleRead), GetRole)
			adminRole.GET("/permissions", middleware.RequirePermission(model.PermAdminPermissionList), GetRolePermissions)
			adminRole.PUT("", middleware.RequirePermission(model.PermAdminRoleUpdate), UpdateRole)
			adminRole.DELETE("", middleware.RequirePermission(model.PermAdminRoleDelete), DeleteRole)
			adminRole.POST("/permissions", middleware.RequirePermission(model.PermAdminPermissionAssign), AssignPermission)
			adminRole.DELETE("/permissions", middleware.RequirePermission(model.PermAdminPermissionRevoke), RevokePermission)
		}

		admin.GET("/groups", middleware.RequirePermission(model.PermAdminGroupList), GetGroups)
		admin.POST("/groups", middleware.RequirePermission(model.PermAdminGroupCreate), CreateGroup)
		adminGroup := admin.Group("/groups/:groupID", middleware.SetGroup)
		{
			adminGroup.GET("", middleware.RequirePermission(model.PermAdminGroupRead), GetGroup)
			adminGroup.GET("/users", middleware.RequirePermission(model.PermAdminUserList), GetGroupUsers)
			adminGroup.PUT("", middleware.RequirePermission(model.PermAdminGroupUpdate), UpdateGroup)
			adminGroup.DELETE("", middleware.RequirePermission(model.PermAdminGroupDelete), DeleteGroup)
			adminGroup.POST("/users", middleware.RequirePermission(model.PermAdminUserAssign), AssignUserToGroup)
			adminGroup.DELETE("/users", middleware.RequirePermission(model.PermAdminUserRevoke), RemoveUserFromGroup)
		}

		admin.GET("/users", middleware.RequirePermission(model.PermAdminUserList), GetUsers)
		admin.POST("/users", middleware.RequirePermission(model.PermAdminUserCreate), CreateUser)
		adminUser := admin.Group("/users/:userID", middleware.SetUser)
		{
			adminUser.GET("", middleware.RequirePermission(model.PermAdminUserRead), GetUser)
			adminUser.PUT("", middleware.RequirePermission(model.PermAdminUserUpdate), UpdateUser)
			adminUser.DELETE("", middleware.RequirePermission(model.PermAdminUserDelete), DeleteUser)
			adminUser.POST("/picture", middleware.RequirePermission(model.PermAdminUserUpdate), UploadPicture("user"))
		}

		admin.GET("/oauth", middleware.RequirePermission(model.PermAdminOauthList), GetOauthProviders)
		admin.POST("/oauth", middleware.RequirePermission(model.PermAdminOauthCreate), CreateOauthProvider)
		adminOauth := admin.Group("/oauth/:oauthID", middleware.SetOauth)
		{
			adminOauth.GET("", middleware.RequirePermission(model.PermAdminOauthRead), GetOauthProvider)
			adminOauth.PUT("", middleware.RequirePermission(model.PermAdminOauthUpdate), UpdateOauthProvider)
			adminOauth.POST("/picture", middleware.RequirePermission(model.PermAdminOauthUpdate), UploadPicture("oauth"))
			adminOauth.DELETE("", middleware.RequirePermission(model.PermAdminOauthDelete), DeleteOauthProvider)
		}

		admin.GET("/email", middleware.RequirePermission(model.PermAdminSMTPList), GetEmails)
		admin.GET("/smtp", middleware.RequirePermission(model.PermAdminSMTPList), GetSmtps)
		admin.POST("/smtp", middleware.RequirePermission(model.PermAdminSMTPCreate), CreateSmtp)
		adminSmtp := admin.Group("/smtp/:smtpID", middleware.SetSmtp)
		{
			adminSmtp.GET("", middleware.RequirePermission(model.PermAdminSMTPRead), GetSmtp)
			adminSmtp.PUT("", middleware.RequirePermission(model.PermAdminSMTPUpdate), UpdateSmtp)
			adminSmtp.DELETE("", middleware.RequirePermission(model.PermAdminSMTPDelete), DeleteSmtp)

			adminSmtp.GET("/email", middleware.RequirePermission(model.PermAdminSMTPList), GetEmails)
		}

		admin.GET("/webhook", middleware.RequirePermission(model.PermAdminWebhookList), GetWebhooks)
		admin.GET("/webhook/events", middleware.RequirePermission(model.PermAdminWebhookList), GetEventTypes)
		admin.GET("/webhook/history", middleware.RequirePermission(model.PermAdminWebhookList), GetWebhookHistory)
		admin.POST("/webhook", middleware.RequirePermission(model.PermAdminWebhookCreate), CreateWebhook)
		adminWebhook := admin.Group("/webhook/:webhookID", middleware.SetWebhook)
		{
			adminWebhook.GET("", middleware.RequirePermission(model.PermAdminWebhookRead), GetWebhook)
			adminWebhook.PUT("", middleware.RequirePermission(model.PermAdminWebhookUpdate), UpdateWebhook)
			adminWebhook.DELETE("", middleware.RequirePermission(model.PermAdminWebhookDelete), DeleteWebhook)

			adminWebhook.GET("/history", middleware.RequirePermission(model.PermAdminWebhookList), GetWebhookHistory)
		}

		admin.GET("/challenges", middleware.RequirePermission(model.PermAdminChallengeList), GetChallenges)
		admin.GET("/challenges/categories", middleware.RequirePermission(model.PermAdminChallengeList), GetChallengeCategories)
		admin.POST("/challenges", middleware.RequirePermission(model.PermAdminChallengeCreate), CreateChallenge)
		adminChallenge := admin.Group("/challenges/:challengeID", middleware.SetChallenge)
		{
			adminChallenge.GET("", middleware.RequirePermission(model.PermAdminChallengeRead), GetChallenge)
			adminChallenge.GET("/download", middleware.RequirePermission(model.PermAdminChallengeRead),
				middleware.SetChallengeFile, DownloadFile(model.DownloadAttachmentEventType),
			)
			adminChallenge.PUT("", middleware.RequirePermission(model.PermAdminChallengeUpdate), UpdateChallenge)
			adminChallenge.DELETE("", middleware.RequirePermission(model.PermAdminChallengeDelete), DeleteChallenge)
			adminChallenge.POST("/upload", middleware.RequirePermission(model.PermAdminChallengeUpdate), UploadChallengeFile)

			adminChallengeTest := adminChallenge.Group("/test", middleware.RequirePermission(model.PermAdminChallengeTest))
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

		admin.GET("/contests", middleware.RequirePermission(model.PermAdminContestList), GetContests)
		admin.POST("/contests", middleware.RequirePermission(model.PermAdminContestCreate), CreateContest)
		adminContest := admin.Group("/contests/:contestID", middleware.SetContest)
		{
			adminContest.GET("", middleware.RequirePermission(model.PermAdminContestRead), GetContest)
			adminContest.PUT("", middleware.RequirePermission(model.PermAdminContestUpdate), UpdateContest)
			adminContest.DELETE("", middleware.RequirePermission(model.PermAdminContestDelete), DeleteContest)
			adminContest.POST("/picture", middleware.RequirePermission(model.PermAdminContestUpdate), UploadPicture("contest"))
			adminContest.GET("/rank", middleware.RequirePermission(model.PermAdminContestRank), GetTeamRanking)
			adminContest.GET("/scoreboard", middleware.RequirePermission(model.PermAdminContestRank), GetScoreboard)
			adminContest.GET("/timeline", middleware.RequirePermission(model.PermAdminContestRank), GetRankTimeline)

			adminContest.GET("/teams", middleware.RequirePermission(model.PermAdminTeamList), GetTeams)
			adminContestTeam := adminContest.Group("/teams/:teamID", middleware.SetTeam)
			{
				adminContestTeam.GET("", middleware.RequirePermission(model.PermAdminTeamRead), GetTeam)
				adminContestTeam.GET("/users", middleware.RequirePermission(model.PermAdminTeamRead), GetTeammates)
				adminContestTeam.PUT("", middleware.RequirePermission(model.PermAdminTeamUpdate), UpdateTeam)
				adminContestTeam.DELETE("", middleware.RequirePermission(model.PermAdminTeamDelete), DeleteTeam)
				adminContestTeam.POST("/kick", middleware.RequirePermission(model.PermAdminTeamUpdate), KickMember)
				adminContestTeam.POST("/picture", middleware.RequirePermission(model.PermAdminTeamUpdate), UploadPicture("team"))

				adminContestTeam.GET("/flags", middleware.RequirePermission(model.PermAdminTeamRead), GetTeamFlags)

				adminContestTeam.GET("/submissions", middleware.RequirePermission(model.PermAdminTeamRead), GetSubmissions)

				adminContestTeam.GET("/victims", middleware.RequirePermission(model.PermAdminTeamRead), GetVictims)
				adminContainer := adminContestTeam.Group("/victims/:victimID", middleware.RequirePermission(model.PermAdminTeamRead), middleware.SetVictim)
				{
					adminTraffic := adminContainer.Group("/traffic")
					{
						adminTraffic.GET("/download", middleware.SetTrafficFile, DownloadFile(model.DownloadTrafficEventType))
						adminTraffic.GET("", GetTraffics)
					}
				}

				adminContestTeam.GET("/writeups",
					middleware.RequirePermission(model.PermAdminTeamWriteupList), GetWriteUPs,
				)
				adminContestTeam.GET("/writeups/:fileID",
					middleware.RequirePermission(model.PermAdminTeamWriteupRead),
					middleware.SetFile(model.WriteupFileType), DownloadFile(model.DownloadWriteUpEventType),
				)
			}

			adminContest.GET("/notices", middleware.RequirePermission(model.PermAdminNoticeList), GetNotices)
			adminContest.POST("/notices", middleware.RequirePermission(model.PermAdminNoticeCreate), CreateNotice)
			adminContestNotice := adminContest.Group("/notices/:noticeID", middleware.SetNotice)
			{
				adminContestNotice.PUT("", middleware.RequirePermission(model.PermAdminNoticeUpdate), UpdateNotice)
				adminContestNotice.DELETE("", middleware.RequirePermission(model.PermAdminNoticeDelete), DeleteNotice)
			}

			adminContest.GET("/cheats", middleware.RequirePermission(model.PermAdminCheatList), GetCheats)
			adminContest.DELETE("/cheats", middleware.RequirePermission(model.PermAdminCheatDelete), DeleteCheat(true))
			adminContest.POST("/cheats", middleware.RequirePermission(model.PermAdminCheatCreate), CheckCheat)
			adminContestCheat := adminContest.Group("/cheats/:cheatID", middleware.SetCheat)
			{
				adminContestCheat.PUT("", middleware.RequirePermission(model.PermAdminCheatUpdate), UpdateCheat)
				adminContestCheat.DELETE("", middleware.RequirePermission(model.PermAdminCheatDelete), DeleteCheat(false))
			}

			adminContest.GET("/challenges", middleware.RequirePermission(model.PermAdminContestChallengeList), GetContestChallenges)
			adminContest.GET("/challenges/others", middleware.RequirePermission(model.PermAdminChallengeList), GetChallengeNotInContest)
			adminContest.GET("/challenges/categories", middleware.RequirePermission(model.PermAdminContestChallengeList), GetContestChallengeCategories)
			adminContest.POST("/challenges", middleware.RequirePermission(model.PermAdminContestChallengeCreate), AddContestChallenge)
			adminContestChallenge := adminContest.Group("/challenges/:challengeID", middleware.SetContestChallenge)
			{
				adminContestChallenge.GET("", middleware.RequirePermission(model.PermAdminContestChallengeRead), GetContestChallenge)
				adminContestChallenge.PUT("", middleware.RequirePermission(model.PermAdminContestChallengeUpdate), UpdateContestChallenge)
				adminContestChallenge.DELETE("", middleware.RequirePermission(model.PermAdminContestChallengeDelete), DeleteContestChallenge)

				//不允许后期创建和删除
				adminContestChallenge.GET("/flags", middleware.RequirePermission(model.PermAdminContestChallengeFlagList), GetContestFlags)
				adminContestFlag := adminContestChallenge.Group("/flags/:flagID", middleware.SetContestFlag)
				{
					adminContestFlag.GET("", middleware.RequirePermission(model.PermAdminContestChallengeFlagRead), GetContestFlag)
					adminContestFlag.PUT("", middleware.RequirePermission(model.PermAdminContestChallengeFlagUpdate), UpdateContestFlag)
				}
			}

			adminContestImages := adminContest.Group("/images", middleware.RequirePermission(model.PermAdminImagePull))
			{
				adminContestImages.GET("", GetContestChallengeImage)
				adminContestImages.POST("", WarmUpContestChallengeImage)
			}

			adminContestVictim := adminContest.Group("/victims", middleware.RequirePermission(model.PermAdminVictimControl))
			{
				adminContestVictim.GET("", GetContestVictims)
				adminContestVictim.POST("", StartContestVictims)
				adminContestVictim.DELETE("", StopContestVictims)
			}
		}

		admin.GET("/files", middleware.RequirePermission(model.PermAdminFileList), GetFiles)
		admin.DELETE("/files", middleware.RequirePermission(model.PermAdminFileDelete), DeleteFiles)
		admin.GET("/files/:fileID", middleware.RequirePermission(model.PermAdminFileRead), middleware.SetFile(""), DownloadFile(model.DownloadFileEventType))

		admin.GET("/logs", middleware.RequirePermission(model.PermAdminLogRead), GetLogs)
	}
	return router
}
