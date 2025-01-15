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
	router.POST("/admin/login", AdminLogin)

	_ = router.Group("/admin")

	return router

}
