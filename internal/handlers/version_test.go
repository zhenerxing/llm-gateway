package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	
	"github.com/gin-gonic/gin"
)

func TestVersion (t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	r := gin.New()
	r.GET("/version",Version)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet,"/version",nil)
	
	r.ServeHTTP(rec,req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d",http.StatusOK,rec.Code)
	}

	type VersionResp struct {
		Code int    `json:"code"`
		Data string `json:"data"`
		Msg  string `json:"msg"`
	}

	var got VersionResp

	if err := json.Unmarshal(rec.Body.Bytes(),&got); err !=nil{
		t.Fatalf("failed to parse json: %v; raw = %s",err , rec.Body.String())
	}
	if got.Code != 200 {
		t.Fatalf("expected code == 200; got %d", got.Code)
	}
	if got.Data != "v1.0.0" {
		t.Fatalf("expected data == v1.0.0; got %q", got.Data)
	}
	if got.Msg == "" {
		t.Fatalf("expected msg not empty; got %q", got.Msg)
	}


}