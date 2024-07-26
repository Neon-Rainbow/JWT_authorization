package main

import (
	"JWT_authorization/config"
	"JWT_authorization/route"
	"JWT_authorization/util"
	"log"
	"sync"
)

func main() {
	err := config.LoadConfig("./config.json")
	if err != nil {
		log.Println("ErrorMessage loading config")
		return
	}

	err = util.Init()
	if err != nil {
		log.Println("ErrorMessage initializing the database")
		return
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		route.StartGRPCServer()
	}()

	go func() {
		defer wg.Done()
		route.StartHTTPServer()
	}()

	wg.Wait()
}
