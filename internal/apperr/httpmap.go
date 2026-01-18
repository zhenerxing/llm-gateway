package apperr

import (
	"net/http"
)

func HTTPStatus(err error) int {
	ae, ok := As(err)
	if !ok {
		return http.StatusInternalServerError
	}

	switch ae.Code {
	case PLATFORM_REQUEST_INVALID:
		return http.StatusBadRequest

	case PLATFORM_AUTH_MISSING_API_KEY, PLATFORM_AUTH_INVALID_API_KEY:
		return http.StatusUnauthorized

	case PLATFORM_AUTH_FORBIDDEN:
		return http.StatusForbidden

	case PLATFORM_DEPENDENCY_UNAVAILABLE:
		return http.StatusServiceUnavailable

	case UPSTREAM_TIMEOUT:
		return http.StatusGatewayTimeout

	case UPSTREAM_UNAVAILABLE:
		return http.StatusServiceUnavailable

	case PLATFORM_INTERNAL_ERROR:
		fallthrough
	default:
		return http.StatusInternalServerError
	}
}