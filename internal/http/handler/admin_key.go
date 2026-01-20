package handler 

import(
	"net/http"

	"github.com/zhenerxing/llm-gateway/internal/auth"

	"github.com/gin-gonic/gin"
	"github.com/zhenerxing/llm-gateway/internal/apperr"
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
		new_err := apperr.New( apperr.PLATFORM_REQUEST_INVALID , apperr.TypePlatform , "input should be in json")
		_ = c.Error(new_err)
			c.Abort()
			return
	}
	/*
	func New(code string , typ Type , message string) *Error{
	return &Error{Code: code, Type: typ, Message: message}
}
	*/
	out, err := h.authSvc.CreateKey(in)
	if err != nil{
		_ = c.Error(err)
			c.Abort()
		return
	}

	c.JSON(http.StatusCreated, out.KeyInfo)
}

func (h *AdminKeysHandler) List(c *gin.Context) {
	recs, err := h.authSvc.ListKeys()
	if err != nil {
		new_err := apperr.New( apperr.PLATFORM_DEPENDENCY_UNAVAILABLE , apperr.TypePlatform , "query store.List error")
		_ = c.Error(new_err)
			c.Abort()
		return
	}
	c.JSON(http.StatusOK, recs)
}