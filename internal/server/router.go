package server

import(
	"github.com/gin-gonic/gin"

	"github.com/zhenerxing/llm-gateway/internal/handlers"

)

func Router() *gin.Engine {
	r:= gin.New()

	r.GET("/version",handlers.Version)
	r.GET("/healthz",handlers.Healthz)

	return r
}