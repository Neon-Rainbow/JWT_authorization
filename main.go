package main

import (
	"JWT_authorization/config"
	"JWT_authorization/route"
	"JWT_authorization/util"
	"fmt"
	"log"
)

func main() {
	err := config.LoadConfig("./config.json")
	if err != nil {
		log.Println("ErrorMessage loading config")
		return
	}

	util.Init()

	r := route.NewRouter()
	
	go route.StartGRPCServer()

	addr := fmt.Sprintf("%v:%v", config.GetConfig().Address, config.GetConfig().Port)
	err = r.Run(addr)
	if err != nil {
		log.Println("ErrorMessage starting server")
		return
	}
}
