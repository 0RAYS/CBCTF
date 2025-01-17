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

	user := router.Group("/user", middleware.CheckLogin)
	user.GET("/info", GetUser)
	user.PUT("/password", ChangePassword)
	user.PUT("/update", UpdateUser)

	admin := router.Group("/admin", middleware.CheckLogin, middleware.CheckAdmin)
	admin.GET("/info", GetAdmin)
	admin.PUT("/password", AdminChangePassword)

	return router

}
