package main

import (
	"ddtd/server"
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

func main() {
	var config Config
	cleanenv.ReadEnv(&config)
	server, err := server.NewServer(config.Server)
	if err != nil {
		log.Fatal(err)
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
