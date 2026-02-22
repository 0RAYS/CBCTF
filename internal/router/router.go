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
		user.GET("", middleware.RBAC(model.PermSelfRead), GetUser)
		user.PUT("/password", middleware.RBAC(model.PermSelfUpdate), ChangePwd)
		user.PUT("", middleware.RBAC(model.PermSelfUpdate), UpdateUser)
		user.DELETE("", middleware.RBAC(model.PermSelfDelete), DeleteUser)
		user.POST("/picture", middleware.RBAC(model.PermSelfUpdate), UploadPicture("self-user"))
		user.POST("/activate",
			middleware.RBAC(model.PermSelfActivate),
			middleware.RateLimit("activate", 1, time.Minute), ActivateEmail,
		)
	}

	contest := auth.Group("/contests/:contestID", middleware.CheckRole(false), middleware.SetContest)
	{
		contest.GET("", middleware.RBAC(model.PermUserContestRead), GetContest)
		contest.GET("/rank", middleware.RBAC(model.PermUserContestRank), GetTeamRanking)
		contest.GET("/scoreboard", middleware.RBAC(model.PermUserContestRank), GetScoreboard)
		contest.GET("/timeline", middleware.RBAC(model.PermUserContestRank), GetRankTimeline)
		contest.POST("/teams/join", middleware.RBAC(model.PermUserTeamJoin),
			middleware.ContestIsNotOver, middleware.CheckVerified, JoinTeam,
		)
		contest.POST("/teams/create", middleware.RBAC(model.PermUserTeamCreate),
			middleware.ContestIsNotOver, middleware.CheckVerified, CreateTeam,
		)

		contestTeam := contest.Group("/teams/me", middleware.CheckVerified, middleware.SetTeamByUser)
		{
			contestTeam.GET("", middleware.RBAC(model.PermUserTeamRead), GetTeam)
			contestTeam.GET("/captcha", middleware.RBAC(model.PermUserTeamRead), GetTeamCaptcha)
			contestTeam.GET("/users", middleware.RBAC(model.PermUserTeamRead), GetTeammates)
			contestTeam.PUT("/captcha", middleware.RBAC(model.PermUserTeamUpdate),
				middleware.ContestIsNotOver, middleware.CheckCaptain, UpdateCaptcha,
			)
			contestTeam.PUT("", middleware.RBAC(model.PermUserTeamUpdate),
				middleware.ContestIsNotOver, middleware.CheckCaptain, UpdateTeam,
			)
			contestTeam.POST("/picture", middleware.RBAC(model.PermUserTeamUpdate),
				middleware.ContestIsNotOver, middleware.CheckCaptain, UploadPicture("team"),
			)
			contestTeam.DELETE("", middleware.RBAC(model.PermUserTeamDelete),
				middleware.ContestIsComing, middleware.CheckCaptain, DeleteTeam,
			)
			contestTeam.POST("/kick", middleware.RBAC(model.PermUserTeamUpdate),
				middleware.ContestIsComing, middleware.CheckCaptain, KickMember,
			)
			contestTeam.POST("/leave", middleware.RBAC(model.PermUserTeamRead), middleware.ContestIsComing, LeaveTeam)
		}

		// 比赛公告
		{
			contest.GET("/notices", middleware.RBAC(model.PermUserNoticeList), GetNotices)
		}

		contest.GET("/challenges", middleware.RBAC(model.PermUserChallengeList),
			middleware.CheckVerified, middleware.SetTeamByUser, middleware.CheckBanned, middleware.ContestIsNotComing, GetContestChallenges,
		)
		contest.GET("/challenges/categories", middleware.RBAC(model.PermUserChallengeList),
			middleware.CheckVerified, middleware.SetTeamByUser, middleware.CheckBanned, middleware.ContestIsNotComing, GetContestChallengeCategories,
		)
		contestChallenge := contest.Group(
			"/challenges/:challengeID",
			middleware.CheckVerified, middleware.SetTeamByUser, middleware.CheckBanned, middleware.ContestIsNotComing, middleware.SetContestChallenge,
		)
		{
			contestChallenge.GET("", middleware.RBAC(model.PermUserChallengeRead), GetContestChallengeStatus)
			contestChallenge.POST("/init",
				middleware.RBAC(model.PermUserChallengeInit),
				middleware.RateLimit("init_flag", 1, time.Minute),
				middleware.ContestIsRunning, middleware.CheckSolved, InitTeamFlag,
			)
			contestChallenge.GET("/attachment",
				middleware.RBAC(model.PermUserChallengeRead),
				middleware.RateLimit("download_attachment", 10, time.Minute),
				middleware.SetAttachmentFile(false), DownloadFile(model.DownloadAttachmentEventType),
			)
			contestChallenge.POST("/reset",
				middleware.RBAC(model.PermUserChallengeReset),
				middleware.RateLimit("init_flag", 1, time.Minute),
				middleware.ContestIsRunning, middleware.CheckIfGenerated, middleware.CheckSolved, ResetTeamFlag,
			)
			contestChallenge.POST("/start",
				middleware.RBAC(model.PermUserVictimControl),
				middleware.CheckChallengeType(model.PodsChallengeType),
				middleware.RateLimit("start_victim", 1, time.Minute), middleware.CheckTeamVictimCount, middleware.CheckIfGenerated, StartVictim,
			)
			contestChallenge.POST("/increase",
				middleware.RBAC(model.PermUserVictimControl),
				middleware.CheckChallengeType(model.PodsChallengeType), middleware.ContestIsRunning,
				middleware.CheckIfGenerated, IncreaseVictimDuration,
			)
			contestChallenge.POST("/stop", middleware.RBAC(model.PermUserVictimControl),
				middleware.CheckChallengeType(model.PodsChallengeType), middleware.CheckIfGenerated, StopVictim,
			)
			contestChallenge.POST("/submit",
				middleware.RBAC(model.PermUserChallengeSubmit),
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
			contestWriteUp.POST("", middleware.RBAC(model.PermUserWriteupUpload), UploadWriteUp)
			contestWriteUp.GET("", middleware.RBAC(model.PermUserWriteupList), GetWriteUPs)
		}
	}

	admin := auth.Group("/admin", middleware.CheckRole(true))
	{
		admin.GET("/ip", middleware.RBAC(model.PermAdminIPSearch), SearchIP)

		admin.GET("/me", GetAdmin)
		admin.PUT("/me/password", AdminChangePassword)
		admin.PUT("/me", UpdateAdmin)
		admin.POST("/me/picture", UploadPicture("admin"))
		admin.POST("", CreateAdmin)

		admin.GET("/models", middleware.RBAC(model.PermAdminModelsSearch), GetAllowQueryModels)
		admin.GET("/search", middleware.RBAC(model.PermAdminModelsSearch), Search)

		// 系统管理
		adminSystem := admin.Group("/system")
		{
			adminSystem.GET("/status", middleware.RBAC(model.PermAdminSystemRead), SystemStatus)
			adminSystem.GET("/config", middleware.RBAC(model.PermAdminSystemRead), SystemConfig)
			adminSystem.PUT("/config", middleware.RBAC(model.PermAdminSystemUpdate), UpdateSystem)
			adminSystem.POST("/restart", middleware.RBAC(model.PermAdminSystemRestart), RestartSystem)
		}

		admin.GET("/permissions", middleware.RBAC(model.PermAdminPermissionList), GetPermissions)
		adminPermission := admin.Group("/permissions/:permissionID", middleware.SetPermission)
		{
			adminPermission.PUT("", middleware.RBAC(model.PermAdminPermissionUpdate), UpdatePermission)
		}

		admin.GET("/roles", middleware.RBAC(model.PermAdminRoleList), GetRoles)
		admin.POST("/roles", middleware.RBAC(model.PermAdminRoleCreate), CreateRole)
		adminRole := admin.Group("/roles/:roleID", middleware.SetRole)
		{
			adminRole.GET("", middleware.RBAC(model.PermAdminRoleRead), GetRole)
			adminRole.GET("/permissions", middleware.RBAC(model.PermAdminPermissionList), GetRolePermissions)
			adminRole.PUT("", middleware.RBAC(model.PermAdminRoleUpdate), UpdateRole)
			adminRole.DELETE("", middleware.RBAC(model.PermAdminRoleDelete), DeleteRole)
			adminRole.POST("/permissions", middleware.RBAC(model.PermAdminRoleAssign), AssignPermission)
			adminRole.DELETE("/permissions", middleware.RBAC(model.PermAdminRoleRevoke), RevokePermission)
		}

		admin.GET("/groups", middleware.RBAC(model.PermAdminGroupList), GetGroups)
		admin.POST("/groups", middleware.RBAC(model.PermAdminGroupCreate), CreateGroup)
		adminGroup := admin.Group("/groups/:groupID", middleware.SetGroup)
		{
			adminGroup.GET("", middleware.RBAC(model.PermAdminGroupRead), GetGroup)
			adminGroup.GET("/users", middleware.RBAC(model.PermAdminUserList), GetGroupUsers)
			adminGroup.PUT("", middleware.RBAC(model.PermAdminGroupUpdate), UpdateGroup)
			adminGroup.DELETE("", middleware.RBAC(model.PermAdminGroupDelete), DeleteGroup)
			adminGroup.POST("/users", middleware.RBAC(model.PermAdminUserAssign), AssignUserToGroup)
			adminGroup.DELETE("/users", middleware.RBAC(model.PermAdminUserRevoke), RemoveUserFromGroup)
		}

		admin.GET("/users", middleware.RBAC(model.PermAdminUserList), GetUsers)
		admin.POST("/users", middleware.RBAC(model.PermAdminUserCreate), CreateUser)
		adminUser := admin.Group("/users/:userID", middleware.SetUser)
		{
			adminUser.GET("", middleware.RBAC(model.PermAdminUserRead), GetUser)
			adminUser.PUT("", middleware.RBAC(model.PermAdminUserUpdate), UpdateUser)
			adminUser.DELETE("", middleware.RBAC(model.PermAdminUserDelete), DeleteUser)
			adminUser.POST("/picture", middleware.RBAC(model.PermAdminUserUpdate), UploadPicture("user"))
		}

		admin.GET("/oauth", middleware.RBAC(model.PermAdminOauthList), GetOauthProviders)
		admin.POST("/oauth", middleware.RBAC(model.PermAdminOauthCreate), CreateOauthProvider)
		adminOauth := admin.Group("/oauth/:oauthID", middleware.SetOauth)
		{
			adminOauth.PUT("", middleware.RBAC(model.PermAdminOauthUpdate), UpdateOauthProvider)
			adminOauth.POST("/picture", middleware.RBAC(model.PermAdminOauthUpdate), UploadPicture("oauth"))
			adminOauth.DELETE("", middleware.RBAC(model.PermAdminOauthDelete), DeleteOauthProvider)
		}

		admin.GET("/email", middleware.RBAC(model.PermAdminSMTPList), GetEmails)
		admin.GET("/smtp", middleware.RBAC(model.PermAdminSMTPList), GetSmtps)
		admin.POST("/smtp", middleware.RBAC(model.PermAdminSMTPCreate), CreateSmtp)
		adminSmtp := admin.Group("/smtp/:smtpID", middleware.SetSmtp)
		{
			adminSmtp.PUT("", middleware.RBAC(model.PermAdminSMTPUpdate), UpdateSmtp)
			adminSmtp.DELETE("", middleware.RBAC(model.PermAdminSMTPDelete), DeleteSmtp)

			adminSmtp.GET("/email", middleware.RBAC(model.PermAdminSMTPList), GetEmails)
		}

		admin.GET("/webhook", middleware.RBAC(model.PermAdminWebhookList), GetWebhooks)
		admin.GET("/webhook/events", middleware.RBAC(model.PermAdminWebhookList), GetEventTypes)
		admin.GET("/webhook/history", middleware.RBAC(model.PermAdminWebhookList), GetWebhookHistory)
		admin.POST("/webhook", middleware.RBAC(model.PermAdminWebhookCreate), CreateWebhook)
		adminWebhook := admin.Group("/webhook/:webhookID", middleware.SetWebhook)
		{
			adminWebhook.PUT("", middleware.RBAC(model.PermAdminWebhookUpdate), UpdateWebhook)
			adminWebhook.DELETE("", middleware.RBAC(model.PermAdminWebhookDelete), DeleteWebhook)

			adminWebhook.GET("/history", middleware.RBAC(model.PermAdminWebhookList), GetWebhookHistory)
		}

		admin.GET("/challenges", middleware.RBAC(model.PermAdminChallengeList), GetChallenges)
		admin.GET("/challenges/categories", middleware.RBAC(model.PermAdminChallengeList), GetChallengeCategories)
		admin.POST("/challenges", middleware.RBAC(model.PermAdminChallengeCreate), CreateChallenge)
		adminChallenge := admin.Group("/challenges/:challengeID", middleware.SetChallenge)
		{
			adminChallenge.GET("/download", middleware.RBAC(model.PermAdminChallengeRead),
				middleware.SetChallengeFile, DownloadFile(model.DownloadAttachmentEventType),
			)
			adminChallenge.PUT("", middleware.RBAC(model.PermAdminChallengeUpdate), UpdateChallenge)
			adminChallenge.DELETE("", middleware.RBAC(model.PermAdminChallengeDelete), DeleteChallenge)
			adminChallenge.POST("/upload", middleware.RBAC(model.PermAdminChallengeUpdate), UploadChallengeFile)

			adminChallengeTest := adminChallenge.Group("/test", middleware.RBAC(model.PermAdminChallengeTest))
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

		admin.GET("/contests", middleware.RBAC(model.PermAdminContestList), GetContests)
		admin.POST("/contests", middleware.RBAC(model.PermAdminContestCreate), CreateContest)
		adminContest := admin.Group("/contests/:contestID", middleware.SetContest)
		{
			adminContest.GET("", middleware.RBAC(model.PermAdminContestRead), GetContest)
			adminContest.PUT("", middleware.RBAC(model.PermAdminContestUpdate), UpdateContest)
			adminContest.DELETE("", middleware.RBAC(model.PermAdminContestDelete), DeleteContest)
			adminContest.POST("/picture", middleware.RBAC(model.PermAdminContestUpdate), UploadPicture("contest"))
			adminContest.GET("/rank", middleware.RBAC(model.PermAdminContestRank), GetTeamRanking)
			adminContest.GET("/scoreboard", middleware.RBAC(model.PermAdminContestRank), GetScoreboard)
			adminContest.GET("/timeline", middleware.RBAC(model.PermAdminContestRank), GetRankTimeline)

			adminContest.GET("/teams", middleware.RBAC(model.PermAdminTeamList), GetTeams)
			adminContestTeam := adminContest.Group("/teams/:teamID", middleware.SetTeam)
			{
				adminContestTeam.GET("", middleware.RBAC(model.PermAdminTeamRead), GetTeam)
				adminContestTeam.GET("/users", middleware.RBAC(model.PermAdminTeamRead), GetTeammates)
				adminContestTeam.PUT("", middleware.RBAC(model.PermAdminTeamUpdate), UpdateTeam)
				adminContestTeam.DELETE("", middleware.RBAC(model.PermAdminTeamDelete), DeleteTeam)
				adminContestTeam.POST("/kick", middleware.RBAC(model.PermAdminTeamUpdate), KickMember)
				adminContestTeam.POST("/picture", middleware.RBAC(model.PermAdminTeamUpdate), UploadPicture("team"))

				adminContestTeam.GET("/flags", middleware.RBAC(model.PermAdminTeamRead), GetTeamFlags)

				adminContestTeam.GET("/submissions", middleware.RBAC(model.PermAdminTeamRead), GetSubmissions)

				adminContestTeam.GET("/victims", middleware.RBAC(model.PermAdminTeamRead), GetVictims)
				adminContainer := adminContestTeam.Group("/victims/:victimID", middleware.RBAC(model.PermAdminTeamRead), middleware.SetVictim)
				{
					adminTraffic := adminContainer.Group("/traffic")
					{
						adminTraffic.GET("/download", middleware.SetTrafficFile, DownloadFile(model.DownloadTrafficEventType))
						adminTraffic.GET("", GetTraffics)
					}
				}

				adminContestTeam.GET("/writeups",
					middleware.RBAC(model.PermAdminTeamWriteupList), GetWriteUPs,
				)
				adminContestTeam.GET("/writeups/:fileID",
					middleware.RBAC(model.PermAdminTeamWriteupRead),
					middleware.SetFile(model.WriteupFileType), DownloadFile(model.DownloadWriteUpEventType),
				)
			}

			adminContest.GET("/notices", middleware.RBAC(model.PermAdminNoticeList), GetNotices)
			adminContest.POST("/notices", middleware.RBAC(model.PermAdminNoticeCreate), CreateNotice)
			adminContestNotice := adminContest.Group("/notices/:noticeID", middleware.SetNotice)
			{
				adminContestNotice.PUT("", middleware.RBAC(model.PermAdminNoticeUpdate), UpdateNotice)
				adminContestNotice.DELETE("", middleware.RBAC(model.PermAdminNoticeDelete), DeleteNotice)
			}

			adminContest.GET("/cheats", middleware.RBAC(model.PermAdminCheatList), GetCheats)
			adminContest.DELETE("/cheats", middleware.RBAC(model.PermAdminCheatDelete), DeleteCheat(true))
			adminContest.POST("/cheats", middleware.RBAC(model.PermAdminCheatCreate), CheckCheat)
			adminContestCheat := adminContest.Group("/cheats/:cheatID", middleware.SetCheat)
			{
				adminContestCheat.PUT("", middleware.RBAC(model.PermAdminCheatUpdate), UpdateCheat)
				adminContestCheat.DELETE("", middleware.RBAC(model.PermAdminCheatDelete), DeleteCheat(false))
			}

			adminContest.GET("/challenges", middleware.RBAC(model.PermAdminContestChallengeList), GetContestChallenges)
			adminContest.GET("/challenges/others", middleware.RBAC(model.PermAdminChallengeList), GetChallengeNotInContest)
			adminContest.GET("/challenges/categories", middleware.RBAC(model.PermAdminContestChallengeList), GetContestChallengeCategories)
			adminContest.POST("/challenges", middleware.RBAC(model.PermAdminContestChallengeCreate), AddContestChallenge)
			adminContestChallenge := adminContest.Group("/challenges/:challengeID", middleware.SetContestChallenge)
			{
				adminContestChallenge.PUT("", middleware.RBAC(model.PermAdminContestChallengeUpdate), UpdateContestChallenge)
				adminContestChallenge.DELETE("", middleware.RBAC(model.PermAdminContestChallengeDelete), DeleteContestChallenge)

				//不允许后期创建和删除
				adminContestChallenge.GET("/flags", middleware.RBAC(model.PermAdminContestChallengeFlagList), GetContestFlags)
				adminContestFlag := adminContestChallenge.Group("/flags/:flagID", middleware.SetContestFlag)
				{
					adminContestFlag.PUT("", middleware.RBAC(model.PermAdminContestChallengeFlagUpdate), UpdateContestFlag)
				}
			}

			adminContestImages := adminContest.Group("/images", middleware.RBAC(model.PermAdminImagePull))
			{
				adminContestImages.GET("", GetContestChallengeImage)
				adminContestImages.POST("", WarmUpContestChallengeImage)
			}

			adminContestVictim := adminContest.Group("/victims", middleware.RBAC(model.PermAdminVictimControl))
			{
				adminContestVictim.GET("", GetContestVictims)
				adminContestVictim.POST("", StartContestVictims)
				adminContestVictim.DELETE("", StopContestVictims)
			}
		}

		admin.GET("/files", middleware.RBAC(model.PermAdminFileList), GetFiles)
		admin.DELETE("/files", middleware.RBAC(model.PermAdminFileDelete), DeleteFiles)
		admin.GET("/files/:fileID", middleware.RBAC(model.PermAdminFileRead), middleware.SetFile(""), DownloadFile(model.DownloadFileEventType))

		admin.GET("/logs", middleware.RBAC(model.PermAdminLogRead), GetLogs)
	}
	return router
}
