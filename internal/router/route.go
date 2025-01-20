package router

import (
	"CBCTF/internal/config"
	"CBCTF/internal/middleware"
	"github.com/gin-gonic/gin"
)

func Init() *gin.Engine {
	gin.SetMode(config.Env.GetString("gin.mode"))
	router := gin.New()
	router.MaxMultipartMemory = config.Env.GetInt64("upload.max") << 20
	router.Use(middleware.Logger(), gin.Recovery(), middleware.Trace, middleware.Cors, middleware.I18n())

	router.POST("/register", Register)
	router.POST("/login", Login)
	router.POST("/admin/login", AdminLogin)

	auth := router.Group("", middleware.CheckLogin)

	user := auth.Group("/user", middleware.CheckRole("user"))
	user.GET("/info", GetUser)
	user.POST("/password", ChangePassword)
	user.POST("/update", UpdateUser)
	user.POST("/delete", DeleteUser)
	user.POST("/avatar")

	contest := auth.Group("/contest", middleware.CheckRole("user"), middleware.CheckVerified)
	contest.GET("/list", GetContests)
	contest.GET("/:contestID/info", middleware.SetContestID, GetContest)
	contest.GET("/:contestID/team/info", middleware.SetContestID, GetTeam)
	contest.GET("/:contestID/team/list", middleware.SetContestID, GetTeams)
	contest.POST("/:contestID/team/join", middleware.SetContestID, JoinTeam)
	contest.POST("/:contestID/team/create", middleware.SetContestID, CreateTeam)
	contest.POST("/:contestID/team/update", middleware.SetContestID, middleware.CheckCaptain, UpdateTeam)
	contest.POST("/:contestID/team/delete", middleware.SetContestID, middleware.CheckCaptain, DeleteTeam)
	contest.POST("/:contestID/team/kick", middleware.SetContestID, middleware.CheckCaptain, KickMember)
	contest.POST("/:contestID/team/leave", middleware.SetContestID, LeaveTeam)

	admin := auth.Group("/admin", middleware.CheckRole("admin"))
	admin.GET("/info", GetAdmin)
	admin.POST("/password", AdminChangePassword)
	admin.POST("/update", UpdateAdmin)

	adminUser := admin.Group("/user")
	adminUser.GET("/list", GetUsers)
	adminUser.POST("/create", CreateUser)
	adminUser.GET("/:userID/info", middleware.SetUserID, GetUser)
	adminUser.POST("/:userID/update", middleware.SetUserID, UpdateUser)
	adminUser.POST("/:userID/delete", middleware.SetUserID, DeleteUser)

	adminContest := admin.Group("/contest")
	adminContest.GET("/list", GetContests)
	adminContest.POST("/create", CreateContest)
	adminContest.GET("/:contestID/info", middleware.SetContestID, GetContest)
	adminContest.POST("/:contestID/update", middleware.SetContestID, UpdateContest)
	adminContest.POST("/:contestID/delete", middleware.SetContestID, DeleteContest)
	adminContest.POST("/:contestID/avatar", middleware.SetContestID)

	adminContestTeam := adminContest.Group("/:contestID/team", middleware.SetContestID)
	adminContestTeam.GET("/list", GetTeams)
	adminContestTeam.GET("/:teamID/info", middleware.SetTeamID, GetTeam)
	adminContestTeam.GET("/:teamID/list", middleware.SetTeamID, GetTeamUsers)
	adminContestTeam.POST("/:teamID/update", middleware.SetTeamID, UpdateTeam)
	adminContestTeam.POST("/:teamID/delete", middleware.SetTeamID, DeleteTeam)
	adminContestTeam.POST("/:teamID/kick", middleware.SetTeamID, KickMember)
	adminContestTeam.POST("/:teamID/avatar", middleware.SetTeamID)

	return router

}
