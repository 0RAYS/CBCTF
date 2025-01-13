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

	router.Use(middleware.Logger(), gin.Recovery(), middleware.Trace, middleware.Cors)

	router.POST("/register", Register)
	router.POST("/login", Login)

	user := router.Group("", middleware.CheckLogin)

	user.POST("/upload", Upload)
	user.GET("/upload/:fileID", Download)

	user.GET("/contests", GetContests)
	user.GET("/teams", GetTeams)
	user.GET("/users", GetUsers)

	admin := router.Group("/admin", middleware.CheckLogin, middleware.CheckAdmin)

	fileAdmin := admin.Group("/files")
	fileAdmin.GET("", GetFiles)
	fileAdmin.DELETE("/:fileID", DeleteFile)

	userAdmin := admin.Group("/users")
	userAdmin.GET("", GetUsers)
	userAdmin.POST("", CreateUser)
	userAdmin.GET("/:userID", middleware.SetUser, GetUser)
	userAdmin.PATCH("/:userID", middleware.SetUser, UpdateUser)
	userAdmin.DELETE("/:userID", middleware.SetUser, DeleteUser)

	teamAdmin := admin.Group("/teams")
	teamAdmin.GET("", GetTeams)
	teamAdmin.POST("", CreateTeam)
	teamAdmin.GET("/:teamID", middleware.SetTeam, GetTeam)
	teamAdmin.PATCH("/:teamID", middleware.SetTeam, UpdateTeam)
	teamAdmin.DELETE("/:teamID", middleware.SetTeam, DeleteTeam)

	contestAdmin := admin.Group("/contests")
	contestAdmin.GET("", GetContests)
	contestAdmin.POST("", CreateContest)
	contestAdmin.GET("/:contestID", middleware.SetContest, GetContest)
	contestAdmin.PATCH("/:contestID", middleware.SetContest, UpdateContest)
	contestAdmin.DELETE("/:contestID", middleware.SetContest, DeleteContest)
	return router
}
