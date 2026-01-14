package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)
type Response struct{
	Code int	     `json:"code"` // 业务码，0或200表示成功，其余表示失败
	Msg  string      `json:"msg"` // 提示信息
	Data interface{} `json:"data"` // 数据，使用 interface{} 以便接收任意类型
}

func Version(c *gin.Context){
	// 返回JSON
	c.JSON(http.StatusOK,Response{
		Code:200,
		Msg:"sucesss",
		Data:"v1.0.0",
	})
}