package middleware

import(
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)
//管理员认证，设在中间层，持有令牌的人才可以入内
func AdminAuthMiddleware()gin.HandlerFunc{
	// 获取令牌
	adminToken := os.Getenv("ADMIN_TOKEN")
	if adminToken == ""{
		adminToken = "dev-admin-token" //先设置一个默认值
	}
	return func(c *gin.Context){
		got := c.GetHeader("X-Admin-Token")
		if got == "" || got != adminToken{
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized_admin",
			})
			return
		}
		c.Next()
	}
}