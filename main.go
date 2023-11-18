package main

import (
	"net/http"

	"github.com/f7ed0/HebergOS-dockerlink/docker"
	"github.com/f7ed0/HebergOS-dockerlink/handling"
	"github.com/f7ed0/HebergOS-dockerlink/logger"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func main() {

	godotenv.Load(".env")

	err := logger.InitDefault("out.log")
	
	if err != nil {
		panic(err)
	}
	
	go docker.FetchStat()

	mux := http.NewServeMux()

	mux.HandleFunc("/v1/container",handling.Container)
	mux.HandleFunc("/v1/container/stats",handling.Stats)
	mux.HandleFunc("/v1/container/start",handling.StartDocker)
	mux.HandleFunc("/v1/container/stop",handling.StopDocker)
	mux.HandleFunc("/v1/container/restart",handling.RestartDocker)

	mux.HandleFunc("/v1/git",handling.Git)
	mux.HandleFunc("/v1/git/head",handling.GitHead)
	mux.HandleFunc("/v1/git/branch",handling.GitBranch)
	mux.HandleFunc("/v1/git/branches",handling.GitBranches)

	mux.HandleFunc("/v1/env",handling.Env)

	logger.Default.Log("INFO","dockerlink started.")
	handler := cors.AllowAll().Handler(mux)
	err = http.ListenAndServe("0.0.0.0:7200",handler)
	if err != nil {
		logger.Default.LogPanic(err.Error())
	}
}

