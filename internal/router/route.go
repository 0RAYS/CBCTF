package router

import (
	"CBCTF/internal/config"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
	"github.com/gin-gonic/gin"
)

func Init() *gin.Engine {
	gin.SetMode(config.Env.GetString("gin.mode"))
	router := gin.New()
	router.MaxMultipartMemory = config.Env.GetInt64("upload.max") << 20
	router.Use(middleware.Logger(), gin.Recovery(), middleware.Trace, middleware.Cors, middleware.I18n())

	router.GET("/download/:fileID", middleware.SetFileID, Download)
	router.POST("/register", Register)
	router.POST("/login", Login)
	router.POST("/admin/login", AdminLogin)

	auth := router.Group("", middleware.CheckLogin)
	auth.POST("/upload", Upload)

	user := auth.Group("/user", middleware.CheckRole("user"))
	user.GET("/info", GetUser)
	user.POST("/password", ChangePassword)
	user.POST("/update", UpdateUser)
	user.POST("/delete", DeleteUser)
	user.POST("/avatar", Avatar(model.User{}))

	contest := auth.Group("/contest", middleware.CheckRole("user"))
	contest.GET("/list", GetContests)
	contest.GET("/:contestID/info", middleware.SetContestID, GetContest)
	contest.GET("/:contestID/team/info", middleware.SetContestID, GetTeam)
	contest.GET("/:contestID/team/captcha", middleware.SetContestID, GetTeamCaptcha)
	contest.GET("/:contestID/team/list", middleware.SetContestID, GetTeammates)
	contest.POST("/:contestID/team/join", middleware.CheckVerified, middleware.SetContestID, JoinTeam)
	contest.POST("/:contestID/team/create", middleware.CheckVerified, middleware.SetContestID, CreateTeam)
	contest.POST("/:contestID/team/update", middleware.CheckVerified, middleware.SetContestID, middleware.CheckCaptain, UpdateTeam)
	contest.POST("/:contestID/team/avatar", middleware.CheckVerified, middleware.SetContestID, middleware.CheckCaptain, Avatar(model.Team{}))
	contest.POST("/:contestID/team/delete", middleware.CheckVerified, middleware.SetContestID, middleware.CheckCaptain, DeleteTeam)
	contest.POST("/:contestID/team/kick", middleware.CheckVerified, middleware.SetContestID, middleware.CheckCaptain, KickMember)
	contest.POST("/:contestID/team/leave", middleware.CheckVerified, middleware.SetContestID, LeaveTeam)

	admin := auth.Group("/admin", middleware.CheckRole("admin"))
	admin.GET("/info", GetAdmin)
	admin.POST("/password", AdminChangePassword)
	admin.POST("/update", UpdateAdmin)
	admin.POST("/avatar", Avatar(model.Admin{}))

	adminUser := admin.Group("/user")
	adminUser.GET("/list", GetUsers)
	adminUser.POST("/create", CreateUser)
	adminUser.GET("/:userID/info", middleware.SetUserID, GetUser)
	adminUser.POST("/:userID/update", middleware.SetUserID, UpdateUser)
	adminUser.POST("/:userID/delete", middleware.SetUserID, DeleteUser)
	adminUser.POST("/:userID/avatar", middleware.SetUserID, Avatar(model.User{}))

	adminContest := admin.Group("/contest")
	adminContest.GET("/list", GetContests)
	adminContest.POST("/create", CreateContest)
	adminContest.GET("/:contestID/info", middleware.SetContestID, GetContest)
	adminContest.POST("/:contestID/update", middleware.SetContestID, UpdateContest)
	adminContest.POST("/:contestID/delete", middleware.SetContestID, DeleteContest)
	adminContest.POST("/:contestID/avatar", middleware.SetContestID, Avatar(model.Contest{}))

	adminContestTeam := adminContest.Group("/:contestID/team", middleware.SetContestID)
	adminContestTeam.GET("/list", GetTeams)
	adminContestTeam.GET("/:teamID/info", middleware.SetTeamID, GetTeam)
	adminContestTeam.GET("/:teamID/list", middleware.SetTeamID, GetTeamUsers)
	adminContestTeam.POST("/:teamID/update", middleware.SetTeamID, UpdateTeam)
	adminContestTeam.POST("/:teamID/delete", middleware.SetTeamID, DeleteTeam)
	adminContestTeam.POST("/:teamID/kick", middleware.SetTeamID, KickMember)
	adminContestTeam.POST("/:teamID/avatar", middleware.SetTeamID, Avatar(model.Team{}))

	return router

}
