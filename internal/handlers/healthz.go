package handlers
import "github.com/gin-gonic/gin"
import "net/http"

func Healthz(c *gin.Context ){
	// 比较简单的是否成功链接的判断
	c.JSON(http.StatusOK,gin.H{"ok":true})
}


		

