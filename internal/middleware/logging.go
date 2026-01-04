package middleware

import (

	"time"
	"go.uber.org/zap"
	"github.com/gin-gonic/gin"
)

func LoggingMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func (c *gin.Context){
		
		start := time.Now()

		//开始计时之后就执行任务
		c.Next()

		// 获取任务的状态
		latency := time.Since(start)
		status := c.Writer.Status()

		//获取的是路由路径，这样更便与查找问题所在
		//但是c.FullPath()在404的时候不会返回路径，所以增加fallback
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}


		ridAny,_ := c.Get("request_id")
		rid, _ := ridAny.(string)

		logger.Info(
			"http access",
			zap.String("request_id", rid),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.Int("status", status),
			zap.Int64("latency_ms", latency.Milliseconds()),
		)
	}
}