package server

import(
	"github.com/gin-gonic/gin"

	"github.com/zhenerxing/llm-gateway/internal/http/handler"

	"github.com/zhenerxing/llm-gateway/internal/http/middleware"

	"github.com/zhenerxing/llm-gateway/internal/auth"

	"github.com/zhenerxing/llm-gateway/internal/audit"

	"go.uber.org/zap"

)

func Router(logger *zap.Logger,store auth.KeyStore, authSvc *auth.Service, auditStore audit.AuditStore) *gin.Engine {
	// 启动引擎，不包含日志和报错
	r:= gin.New()
	
	r.Use(middleware.RequestIDMiddleware())
	r.Use(middleware.LoggingMiddleware(logger))
	r.Use(middleware.ErrorHandler(logger))
	// r.Use(gin.Recovery()) //增加了errorhandler，避免 panic 被它提前处理

	adminKeys := handler.PointerAdminKeysHandler(authSvc)

	admin := r.Group("/admin")
	admin.Use(middleware.AdminAuthMiddleware())
	{
		admin.POST("/keys", adminKeys.Create)
		admin.GET("/keys", adminKeys.List)
		auditH := &handler.AuditHandler{Store: auditStore}
		admin.GET("/audit", auditH.Get)
	}
		

	// public
	r.GET("/version", handler.Version)
	r.GET("/healthz", handler.Healthz)

	// 业务 API：这里保证 Auth 在 Audit 之前
	api := r.Group("/")
	api.Use(auth.AuthMiddleware(store))
	api.Use(middleware.AuditMiddleware(auditStore))
	{
		api.POST("/chat", handler.Chat)
	}
	return r
}