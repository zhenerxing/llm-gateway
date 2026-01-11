package handler 

import(
	"net/http"

	"github.com/zhenerxing/llm-gateway/internal/auth"

	"github.com/gin-gonic/gin"
)

// 使用这里的函数之前必须先使用middleware里面的admin_auth
type AdminKeysHandler struct{
	authSvc *auth.Service
}
func PointerAdminKeysHandler(authSvc *auth.Service) *AdminKeysHandler{
	return &AdminKeysHandler{authSvc:authSvc}
}

func (h *AdminKeysHandler) Create (c *gin.Context){
	var in auth.CreateKeyInput
	if err := c.ShouldBindJSON(&in); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad_request"})
		return
	}

	out, err := h.authSvc.CreateKey(in)
	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_error"})
		return
	}

	c.JSON(http.StatusCreated, out.KeyInfo)
}

func (h *AdminKeysHandler) List(c *gin.Context) {
	recs, err := h.authSvc.ListKeys()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_error"})
		return
	}
	c.JSON(http.StatusOK, recs)
}