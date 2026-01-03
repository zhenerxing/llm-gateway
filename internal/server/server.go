package server

import(
	"log"
	"os"

)

func Run(){
	// 接收http

	// 将http的具体信息交给router去解析分发
	r:= Router()

	// 解析地址，如果没有给出地址，则默认8080
	addr := ":8080"
	if v := os.Getenv("ADDR"); v != ""{
		addr = v
	}

	// 监听端口，如果返回错误则报错
	log.Printf("listening on %s",addr)
	if err := r.Run(addr); err != nil{
		log.Fatal(err)
	}

}