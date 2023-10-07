package handling

import (
	"herbergOS/docker"
	"log"
	"net/http"

	"github.com/docker/docker/api/types"
)

func StartDocker(resp http.ResponseWriter,req *http.Request) {
	if req.Method != "POST" {
		resp.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	qmap := req.URL.Query()

	id,ok := qmap["id"]
	if(!ok) {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	dk,err := docker.NewDockerHandler()
	if err != nil {
		log.Default().Println(err.Error())
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = dk.Client.ContainerStart(dk.Context,id[0],types.ContainerStartOptions{})
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