package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"github.com/gin-gonic/gin"
)

const (
	RequestIDHeader = "X-Request-ID"
	RequestIDKey = "request_id"
)

func RequestIDMiddleware() gin.HandlerFunc {
	return func (c *gin.Context){

		//获取上游来的rid
		rid := c.GetHeader(RequestIDHeader)

		if rid == ""{
			rid = newRequestID()
		}

		//将rid传给客户端
		c.Header(RequestIDHeader,rid)


		//将rid存在当前gin.Context的作用域中
		c.Set(RequestIDKey,rid)

		//让请求继续
		c.Next()
	}

}

func newRequestID() string {
	var b [16]byte
	_, _ = rand.Read(b[:]) 
	return hex.EncodeToString(b[:])
}