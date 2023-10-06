package main

import (
	"herbergOS/docker"
	"herbergOS/handling"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func main() {
	godotenv.Load(".env")
	
	log.Default().Println("Starting...")
	go docker.FetchStat()

	mux := http.NewServeMux()

	mux.HandleFunc("/v1/container/",handling.Container)
	mux.HandleFunc("/v1/container/stats",handling.Stats)
	mux.HandleFunc("/v1/container/start",handling.StartDocker)
	mux.HandleFunc("/v1/container/stop",handling.StopDocker)

	mux.HandleFunc("/v1/git",handling.Git)
	mux.HandleFunc("/v1/git/head",handling.GitHead)
	mux.HandleFunc("/v1/git/branch",handling.GitBranch)
	mux.HandleFunc("/v1/git/branches",handling.GitBranches)

	mux.HandleFunc("/v1/env",handling.Env)

	log.Default().Println("Started !")
	handler := cors.AllowAll().Handler(mux)
	err := http.ListenAndServe("localhost:7200",handler)
	if err != nil {
		log.Default().Fatal(err.Error())
	}
	
}

