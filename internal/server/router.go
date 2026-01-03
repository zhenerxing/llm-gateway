package server

import(
	"github.com/gin-gonic/gin"

	"github.com/zhenerxing/llm-gateway/internal/handlers"

)

func Router() *gin.Engine {
	// 启动引擎，不包含日志和报错
	r:= gin.New()

	// 将具体的任务分发给handlers
	r.GET("/version",handlers.Version)
	r.GET("/healthz",handlers.Healthz)
	
	return r
}