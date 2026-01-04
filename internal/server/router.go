package server

import(
	"github.com/gin-gonic/gin"

	"github.com/zhenerxing/llm-gateway/internal/handlers"

	"github.com/zhenerxing/llm-gateway/internal/middleware"

	"go.uber.org/zap"

)

func Router(logger *zap.Logger) *gin.Engine {
	// 启动引擎，不包含日志和报错
	r:= gin.New()
	
	r.Use(middleware.RequestIDMiddleware())
	r.Use(gin.Recovery())
	r.Use(middleware.LoggingMiddleware(logger))

	// 将具体的任务分发给handlers
	r.GET("/version",handlers.Version)
	r.GET("/healthz",handlers.Healthz)
	
	return r
}