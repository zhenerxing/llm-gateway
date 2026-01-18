package response
/*
{
  "error": {
    "code": "PLATFORM_AUTH_INVALID_API_KEY",
    "message": "invalid api key",
    "type": "platform",
    "request_id": "..."
  }
}
*/
import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	RequestIDHeader = "X-Request-ID"
	RequestIDKey = "request_id"
)

type ErrorBody struct{
	Code      string          `json:"code"`
	Message   string          `json:"message"`
	Type      string          `json:"type"` // "platform" | "upstream"
	RequestID string          `json:"request_id"`
	Details   map[string]any  `json:"details"`
}

// 响应包装，便于后续拓展
type ErrorResponse struct{
	Error ErrorBody `json:"error"`
}

type SuccessResponse struct{
	Data any `json:"data,omitempty"`
	RequestID string `json:"request_id,omitempty"`
}

func OK(c *gin.Context, data any){
	c.JSON(http.StatusOK, SuccessResponse{
		Data:      data,
		RequestID: RequestID(c),
	})
}

func Created(c *gin.Context, data any) {
	c.JSON(http.StatusCreated, SuccessResponse{
		Data:      data,
		RequestID: RequestID(c),
	})
}

func RequestID(c *gin.Context) string {
	// 兼容：RequestIDMiddleware 放 header
	if v := c.GetHeader(RequestIDHeader); v != "" {
		return v
	}
	// 兼容：RequestIDMiddleware 放 context
	if v, ok := c.Get(RequestIDKey); ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}