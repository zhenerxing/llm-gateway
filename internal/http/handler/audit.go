package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zhenerxing/llm-gateway/internal/audit"
)

// 面向接口的编程，这样以后就算数据层变了，也不影响执行层的操作
type AuditHandler struct{
	Store audit.AuditStore
}

// 现有的业务函数，查询
func (h *AuditHandler) Get(c *gin.Context){
	tenantID := c.Query("tenant_id")
	if tenantID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing tenant_id"})
		return
	}

	from := c.Query("from")
	to := c.Query("to")

	limitStr := c.DefaultQuery("limit", "100")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 || limit > 1000 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit (1~1000)"})
		return
	}

	// Query(ctx, tenantID, from, to string, limit int) ([]Record, error)
	rows, err := h.Store.Query(c.Request.Context(), tenantID, from, to, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, rows)
}