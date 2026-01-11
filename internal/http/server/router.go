package server

import(
	"github.com/gin-gonic/gin"

	"github.com/zhenerxing/llm-gateway/internal/http/handler"

	"github.com/zhenerxing/llm-gateway/internal/http/middleware"

	"github.com/zhenerxing/llm-gateway/internal/auth"

	"go.uber.org/zap"

)

func Router(logger *zap.Logger,store auth.KeyStore, authSvc *auth.Service) *gin.Engine {
	// 启动引擎，不包含日志和报错
	r:= gin.New()
	
	r.Use(middleware.RequestIDMiddleware())
	r.Use(gin.Recovery())
	r.Use(middleware.LoggingMiddleware(logger))

	adminKeys := handler.PointerAdminKeysHandler(authSvc)

	admin := r.Group("/admin")
	admin.Use(middleware.AdminAuthMiddleware())
	{
		admin.POST("/keys", adminKeys.Create)
		admin.GET("/keys", adminKeys.List)
	}

	// 将具体的任务分发给handlers
	r.GET("/version",handler.Version)
	r.GET("/healthz",handler.Healthz)

	r.POST("/chat",auth.AuthMiddleware(store),handler.Chat)
	
	return r
}