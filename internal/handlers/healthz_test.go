package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	
	"github.com/gin-gonic/gin"
)

func TestHealthz(t *testing.T){
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("healthz",Healthz)

	req := httptest.NewRequest(http.MethodGet,"/healthz",nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec,req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected statue %d, got %d",http.StatusOK,rec.Code)
	}

	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(),&body); err !=nil{
		t.Fatalf("failed to parse json: %v; raw = %s",err , rec.Body.String())
	}

	ok , exists := body["ok"]
	if !exists{
		t.Fatalf(`expected field "ok" to exist; got %v`,body)
	}

	if ok != true {
		t.Fatalf(`expected "ok" == true; got %v`, ok)
	}

}