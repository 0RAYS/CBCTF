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

	user := router.Group("/user", middleware.CheckLogin, middleware.CheckType("user"))
	user.GET("/info", GetUser)
	user.POST("/password", ChangePassword)
	user.POST("/update", UpdateSelf)

	contest := router.Group("/contest", middleware.CheckLogin)
	contest.GET("/list", GetContests)

	admin := router.Group("/admin", middleware.CheckLogin, middleware.CheckType("admin"))
	admin.GET("/info", GetAdmin)
	admin.POST("/password", AdminChangePassword)
	admin.POST("/update", UpdateAdmin)

	adminUser := admin.Group("/user")
	adminUser.GET("/list", GetUsers)
	adminUser.POST("/:userID/update", UpdateUser)

	adminContest := admin.Group("/contest")
	adminContest.GET("/list", GetContests)
	//adminContest.POST("/create", CreateContest)
	//adminContest.PUT("/update", UpdateContest)
	//adminContest.DELETE("/delete", DeleteContest)

	return router

}
