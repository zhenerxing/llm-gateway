package main

import(
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main(){
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
}