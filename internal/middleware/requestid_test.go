package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/zhenerxing/llm-gateway/internal/handlers"
	
	"github.com/gin-gonic/gin"
)


func TestRequestID(t *testing.T){
	gin.SetMode(gin.TestMode)

	r := gin.New()

	r.Use(RequestIDMiddleware())
	r.GET("healthz",handlers.Healthz)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet,"/healthz",nil)

	r.ServeHTTP(rec,req)

	rid := rec.Header().Get(RequestIDHeader)

	if rid == "" {
		t.Fatalf("expected X-Request-ID header to be set")
	}


}
