package handler

import(
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zhenerxing/llm-gateway/internal/auth"
)

type ChatRequest struct{
	Message string `json:"message"`
}

func Chat(c *gin.Context){
	tenantID,_ := c.Get(auth.CtxKeyTenantID)

	var req ChatRequest
	_ = c.ShouldBindJSON(&req)

	c.JSON(http.StatusOK, gin.H{
		"ok":        true,
		"tenant_id": tenantID,
		"echo":      req.Message,
	})
}