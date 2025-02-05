package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
)

func getDBPoolStats(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	stats := sqlDB.Stats()

	// 打印连接池状态
	fmt.Printf("MaxOpenConnections: %d ", stats.MaxOpenConnections)
	fmt.Printf("OpenConnections: %d ", stats.OpenConnections)
	fmt.Printf("InUse: %d ", stats.InUse)
	fmt.Printf("Idle: %d\n", stats.Idle)
}

// PoolStatsMiddleware 自定义 Gin 中间件，打印每个请求的连接池信息
func PoolStatsMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		getDBPoolStats(db)
		c.Next() // 继续执行请求处理
	}
}
