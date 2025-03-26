package router

import (
	"CBCTF/internel/config"
	"CBCTF/internel/log"
	"CBCTF/internel/middleware"
	"github.com/gin-gonic/gin"
)

func Init() *gin.Engine {
	gin.SetMode(config.Env.Gin.Mode)
	router := gin.New()
	log.Logger.Infof("Trust proxies: %v", config.Env.Gin.Proxies)
	if err := router.SetTrustedProxies(config.Env.Gin.Proxies); err != nil {
		log.Logger.Warningf("Set trusted proxies failed: %v", err)
	}
	router.MaxMultipartMemory = int64(config.Env.Gin.Upload.Max << 20)
	router.Use(
		middleware.Logger(), gin.Recovery(), middleware.SetTrace, middleware.Cors,
		middleware.I18n(), middleware.AccessLog, middleware.RateLimit(), middleware.SetMagic,
	)

	{
		router.POST("/register", Register)
		router.POST("/login", Login)
		router.POST("/admin/login", AdminLogin)
		router.POST("/verify", VerifyEmail)
	}

	auth := router.Group("", middleware.CheckAuth)

	user := auth.Group("user", middleware.CheckRole("user"))
	{
		user.GET("", GetUser)
		user.PUT("/password", ChangePwd)
		user.PUT("", UpdateUser)
		user.DELETE("", DeleteUser)
		user.PUT("/avatar", UploadFile("self-user", "avatar"))
		user.POST("/activate", ActivateEmail)
	}

	contest := auth.Group("contest/:contestID", middleware.CheckRole("user"), middleware.SetContest)
	{
		contest.GET("", GetContest)
		contest.GET("/rank", GetTeamRanking)
	}

	return router
}
