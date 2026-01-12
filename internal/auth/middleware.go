package auth

import(
	"net/http"
	"strings"
	"time"

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
func AuthMiddleware(store KeyStore) gin.HandlerFunc{
	return func (c *gin.Context){
		apiKey := strings.TrimSpace(c.GetHeader(HeaderAPIKey))
		if apiKey == ""{
			c.AbortWithStatusJSON(http.StatusUnauthorized,gin.H{
				"error": "missing api key",
			})
			return
		}
		ki,ok := store.Get(apiKey)
		if !ok || !ki.Active {
			c.AbortWithStatusJSON(http.StatusUnauthorized,gin.H{
				"error": "invalid api key",
			})
			return
		}
		if ki.ExpiresAt != nil && ki.ExpiresAt.Before(time.Now()) {
			c.AbortWithStatusJSON(http.StatusUnauthorized,gin.H{
				"error": "api key expired",
			})
			return
		}
		c.Set(CtxKeyAPIKey,apiKey)
		c.Set(CtxKeyTenantID,ki.TenantID)
		c.Next()

	}
}