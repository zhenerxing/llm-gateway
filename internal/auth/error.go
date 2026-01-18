package auth
// 这个包的作用是定义业务层可识别的错误

import "github.com/zhenerxing/llm-gateway/internal/apperr"

var (
	ErrMissingAPIKey = apperr.New(
		apperr.PLATFORM_AUTH_MISSING_API_KEY,
		apperr.TypePlatform,
		"missing api key",
	)

	ErrInvalidAPIKey = apperr.New(
		apperr.PLATFORM_AUTH_INVALID_API_KEY,
		apperr.TypePlatform,
		"invalid api key",
	)
)