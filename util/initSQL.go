package util

import (
	"JWT_authorization/util/MySQL"
	"JWT_authorization/util/Redis"
	"context"
	"log"
	"time"
)

func Init() {
	// This is a placeholder for the init function

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resultChannel := make(chan bool, 2)
	errorChannel := make(chan error)

	go func() {
		err := MySQL.InitMySQL()
		if err != nil {
			errorChannel <- err
			return
		}
		resultChannel <- true
		return
	}()

	go func() {
		err := Redis.InitRedis()
		if err != nil {
			errorChannel <- err
			return
		}
		resultChannel <- true
		return
	}()

	for i := 0; i < 2; i++ {
		select {
		case <-resultChannel:
			// Successfully initialized the database
		case err := <-errorChannel:
			log.Fatalf("Error initializing the database: %v", err)
		case <-ctx.Done():
			log.Fatalf("Timeout initializing the database")
		}
	}
	return
}
