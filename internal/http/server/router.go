package server

import(
	"github.com/gin-gonic/gin"

	"github.com/zhenerxing/llm-gateway/internal/http/handler"

	"github.com/zhenerxing/llm-gateway/internal/http/middleware"

	"github.com/zhenerxing/llm-gateway/internal/auth"

	"go.uber.org/zap"

)

func Router(logger *zap.Logger,store auth.KeyStore) *gin.Engine {
	// 启动引擎，不包含日志和报错
	r:= gin.New()
	
	r.Use(middleware.RequestIDMiddleware())
	r.Use(gin.Recovery())
	r.Use(middleware.LoggingMiddleware(logger))

	// 将具体的任务分发给handlers
	r.GET("/version",handler.Version)
	r.GET("/healthz",handler.Healthz)

	r.POST("/chat",auth.AuthMiddleware(store),handler.Chat)
	
	return r
}