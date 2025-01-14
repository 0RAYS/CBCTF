package route

import (
	"CBCTF/internal/config"
	"CBCTF/internal/middleware"
	"github.com/gin-gonic/gin"
)

func Init() {
	gin.SetMode(config.Env.GetString("gin.mode"))
	router := gin.New()
	router.MaxMultipartMemory = config.Env.GetInt64("upload.max") << 20
	router.Use(middleware.Logger(), gin.Recovery(), middleware.Trace, middleware.Cors)

}
