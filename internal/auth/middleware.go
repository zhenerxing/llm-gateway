package auth

import(
	"strings"
	"time"
	"fmt"

	"github.com/gin-gonic/gin"
)

const (
	HeaderAPIKey = "X-API-Key"
	CtxKeyTenantID = "tenant_id"
	CtxKeyAPIKey = "api_key"
)

// Middleware 校验 X-API-Key
// - 缺少： 401
// - 不存在/不激活/过期：401
// - 通过：把tenant_id 等信息放到context 里
func AuthMiddleware(store KeyStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Printf("AuthMiddleware hit, X-API-Key=%q\n", c.GetHeader("X-API-Key"))
		apiKey := strings.TrimSpace(c.GetHeader(HeaderAPIKey))
		if apiKey == "" {
			_ = c.Error(ErrMissingAPIKey)
			c.Abort()
			return
		}

		ki, ok := store.Get(apiKey)
		if !ok || !ki.Active {
			_ = c.Error(ErrInvalidAPIKey)
			c.Abort()
			return
		}

		if ki.ExpiresAt != nil && ki.ExpiresAt.Before(time.Now()) {
			// 过期也按 invalid 处理(暂时)
			_ = c.Error(ErrInvalidAPIKey)
			c.Abort()
			return
		}

		c.Set(CtxKeyAPIKey, apiKey)
		c.Set(CtxKeyTenantID, ki.TenantID)
		c.Next()
	}
}