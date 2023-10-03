package handling

import (
	"encoding/json"
	"herbergOS/docker"
	"log"
	"net/http"

	"github.com/docker/docker/api/types"
)


func ContainerList(resp http.ResponseWriter,req *http.Request) {
	if(req.Method != "GET") {
		resp.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	dk,err := docker.NewDockerHandler()
	if err != nil {
		log.Default().Printf("ERR : %v\n",err.Error())
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	containers,err := dk.Client.ContainerList(dk.Context,types.ContainerListOptions{
		All:true,
	})

	if err != nil {
		log.Default().Printf("ERR : %v\n",err.Error())
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	result := map[string]map[string]any{}

	for _,container := range containers {
		result[container.ID] = map[string]any{"names":container.Names,"state":container.State}
	}

	resp.Header().Set("Content-Type", "application/json")

	j := json.NewEncoder(resp)

	j.Encode(result)
	
	return
}

