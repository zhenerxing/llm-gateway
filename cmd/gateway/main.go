package main

import "github.com/zhenerxing/llm-gateway/internal/server"


func main(){
	// 只保留启动入口i
	server.Run()
	/*
	r := gin.New()

	r.GET("/healthz",func(c *gin.Context){
		c.JSON(200,gin.H{"ok":true})
	})

	addr := ":8080"
	if v := os.Getenv("ADDR"); v != ""{
		addr = v
	}

	log.Printf("listening on %s",addr)
	if err := r.Run(addr); err != nil{
		log.Fatal(err)
	}
	*/
}