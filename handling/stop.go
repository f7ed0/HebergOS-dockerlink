package handling

import (
	"herbergOS/docker"
	"log"
	"net/http"

	"github.com/docker/docker/api/types/container"
)

func StopDocker(resp http.ResponseWriter,req *http.Request) {
	if req.Method != "POST" {
		resp.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	Qmap := req.URL.Query()

	id,ok := Qmap["id"]
	if(!ok) {
		resp.WriteHeader(http.StatusBadRequest)
	}

	dk,err := docker.NewDockerHandler()
	if err != nil {
		log.Default().Println(err.Error())
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = dk.Client.ContainerStop(dk.Context,id[0],container.StopOptions{})
	
	if err != nil {
		log.Default().Println(err.Error())
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Header().Set("Content-Type", "text/plain")
		resp.Write([]byte(err.Error()))
		return
	}

	resp.WriteHeader(http.StatusNoContent)
	return
}