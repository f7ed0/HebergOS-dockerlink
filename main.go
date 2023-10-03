package main

import (
	"herbergOS/docker"
	"herbergOS/handling"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")

	log.Default().Println("Starting...")
	go docker.FetchStat()

	http.HandleFunc("/v1/container/list",handling.ContainerList)
	http.HandleFunc("/v1/container/stats",handling.Stats)
	http.HandleFunc("/v1/container/start",handling.StartDocker)
	http.HandleFunc("/v1/container/stop",handling.StopDocker)

	log.Default().Println("Started !")
	err := http.ListenAndServe("localhost:7200",nil)
	if err != nil {
		log.Default().Fatal(err.Error())
	}

	
}

