package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zhenerxing/llm-gateway/internal/audit"
)

const (
	CtxKeyTenantID = "tenant_id"
)

func AuditMiddleware(store audit.AuditStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// 先让后续 handler 执行完
		c.Next()

		// tenant_id：安全获取
		v, ok := c.Get(CtxKeyTenantID)
		tenantID, _ := v.(string)
		if !ok || tenantID == "" {
			return
		}

		// request_id：优先从 context 取（由 RequestIDMiddleware 写入）
		ridAny, _ := c.Get(RequestIDKey) // "request_id"
		rid, _ := ridAny.(string)
		if rid == "" {
			// 兜底从 header 取
			rid = c.GetHeader(RequestIDHeader) // "X-Request-ID"
			if rid == "" {
				rid = "missing"
			}
		}

		endpoint := c.FullPath()
		if endpoint == "" {
			endpoint = c.Request.URL.Path
		}

		rec := audit.AuditInfo{
			RequestID: rid,
			KeyID:     "", // 暂时保持空
			TenantID:  tenantID,
			Endpoint:  endpoint,
			Status:    c.Writer.Status(),
			LatencyMS: time.Since(start).Milliseconds(),
			CreatedAt: time.Now().UTC().Format(time.RFC3339Nano),
		}

		if err := store.Insert(c.Request.Context(), rec); err != nil {
			log.Printf("audit insert failed: %v", err)
		}
	}
}