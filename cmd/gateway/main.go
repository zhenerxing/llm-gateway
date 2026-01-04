package main

import (
	"github.com/zhenerxing/llm-gateway/internal/server"
	"log"
)
func main(){
	// 启动入口i

	//将退出命令放到main中，这样在server中增加了defer，以保存异常退出的日志
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}