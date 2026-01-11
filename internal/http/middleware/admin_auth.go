package middleware

import(
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)
func AdminAuthMiddleware()gin.HandlerFunc{
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