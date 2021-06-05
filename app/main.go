package main

import (
	"cuboid-challenge/app/config"
	"cuboid-challenge/app/db"
	"cuboid-challenge/app/router"
	"fmt"
	"log"
)

func main() {
	config.Load()
	db.Connect()

	r := router.Setup()
	addr := fmt.Sprintf(":%s", config.ENV.Port)

	if err := r.Run(addr); err != nil {
		log.Fatalln("Failed to start the application.", err)
	}
}
