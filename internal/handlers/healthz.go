package handlers
import "github.com/gin-gonic/gin"
import "net/http"

func Healthz(c *gin.Context ){
	c.JSON(http.StatusOK,gin.H{"ok":true})
}


		

