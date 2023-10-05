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

	http.HandleFunc("/v1/container/",handling.Container)
	http.HandleFunc("/v1/container/stats",handling.Stats)
	http.HandleFunc("/v1/container/start",handling.StartDocker)
	http.HandleFunc("/v1/container/stop",handling.StopDocker)

	http.HandleFunc("/v1/git",handling.Git)
	http.HandleFunc("/v1/git/head",handling.GitHead)
	http.HandleFunc("/v1/git/branch",handling.GitBranch)
	http.HandleFunc("/v1/git/branches",handling.GitBranches)


	log.Default().Println("Started !")
	err := http.ListenAndServe("localhost:7200",nil)
	if err != nil {
		log.Default().Fatal(err.Error())
	}
	
}

