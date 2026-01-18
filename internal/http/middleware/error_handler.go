package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/zhenerxing/llm-gateway/internal/apperr"
	"github.com/zhenerxing/llm-gateway/pkg/response"
)

// 下面要完成一个中间层函数
// 作用有两个：
// 一是将panic转化为可控的，error响应给客户端
// 二是将正常传入的error响应给客户端
// 停止当前http请求
func ErrorHandler(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				if logger != nil {
					logger.Error("panic recovered",
						zap.Any("panic", r),
						zap.String("request_id", response.RequestID(c)),
					)
				}
				// panic 一定要回包（如果还没写）
				if !c.Writer.Written() {
					writeError(c, apperr.New(apperr.PLATFORM_INTERNAL_ERROR, apperr.TypePlatform, "internal error"))
				}
				c.Abort()
			}
		}()

		// 关键：先让后续中间件和 handler 执行
		c.Next()

		// 如果已经有响应，就不要覆盖（比如 AuthMiddleware 已经 AbortWithStatusJSON）
		if c.Writer.Written() {
			return
		}
		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err
		writeError(c, err)
		c.Abort()
	}
}

func writeError(c *gin.Context, err error) {
	status := apperr.HTTPStatus(err)

	ae, ok := apperr.As(err)
	if !ok {
		ae = apperr.Wrap(apperr.PLATFORM_INTERNAL_ERROR, apperr.TypePlatform, "internal error", err)
		status = http.StatusInternalServerError
	}

	body := response.ErrorResponse{
		Error: response.ErrorBody{
			Code:      ae.Code,
			Message:   ae.Message,
			Type:      string(ae.Type),
			RequestID: response.RequestID(c),
			Details:   ae.Details,
		},
	}

	if ae.RetryAfter > 0 {
		c.Header("Retry-After", fmt.Sprintf("%d", ae.RetryAfter))
	}

	c.JSON(status, body)
}