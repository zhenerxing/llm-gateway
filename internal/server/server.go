package server

import(
	"log"
	"os"

)

func Run(){

	r:= Router()

	addr := ":8080"
	if v := os.Getenv("ADDR"); v != ""{
		addr = v
	}

	log.Printf("listening on %s",addr)
	if err := r.Run(addr); err != nil{
		log.Fatal(err)
	}

}