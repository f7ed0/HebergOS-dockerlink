package handling

import (
	"encoding/json"
	"herbergOS/docker"
	"log"
	"net/http"
	"time"

	"github.com/docker/docker/api/types"
)


func Container(resp http.ResponseWriter,req *http.Request) {
	switch(req.Method) {
	case http.MethodGet:
		ContainerGet(resp,req)
		return
	case "PUT":
		resp.WriteHeader(http.StatusNotImplemented)
		return
	}
	
	resp.WriteHeader(http.StatusMethodNotAllowed)
	return
	
}

func ContainerGet(resp http.ResponseWriter,req *http.Request) {
	qmap := req.URL.Query()
	j := json.NewEncoder(resp)

	ids,ok := qmap["id"]

	dk,err := docker.NewDockerHandler()
	if err != nil {
		log.Default().Printf("ERR : %v\n",err.Error())
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	running := false
	has_run := false

	if !ok {
		

		r,ok := qmap["running"]

		if ok {
			running = r[0] == "true"
			has_run = true
		}
		
		containers,err := dk.Client.ContainerList(dk.Context,types.ContainerListOptions{
			All: true,
		})
	
		if err != nil {
			log.Default().Printf("ERR : %v\n",err.Error())
			resp.WriteHeader(http.StatusInternalServerError)
			return
		}
	
		ids = []string{}
	
		for _,container := range containers {
			ids = append(ids, container.ID)
		}
	}

	log.Default().Println(ids)

	result := map[string]map[string]any{}
	for _,id := range ids {
		json,err := dk.Client.ContainerInspect(dk.Context,id)
		if err != nil {
			log.Default().Printf("ERR : %v\n",err.Error())
			resp.WriteHeader(http.StatusInternalServerError)
			return
		}
		
		if  has_run && json.State.Running != running {
			continue
		}

		result[id] = map[string]any{"name":json.Name,"state":json.State.Status}

		if json.State.Running {
			t,err := time.Parse(time.RFC3339Nano,json.State.StartedAt)
			if err != nil {
				log.Default().Printf("ERR : %v\n",err.Error())
				resp.WriteHeader(http.StatusInternalServerError)
				return
			}

			result[id]["started_at"] = t.Unix()
		} else {
			result[id]["exit_code"] = json.State.ExitCode
		}
		
	}

	resp.Header().Set("Content-Type", "application/json")
	j.Encode(result)

	return
}

// TODO aled
/*
func ContainerPut(resp http.ResponseWriter,req *http.Request) {
	//qmap := req.URL.Query()
	//j := json.NewEncoder(resp)

	// TODO retrieve data from qmap

	dk,err := docker.NewDockerHandler()
	if err != nil {
		log.Default().Printf("ERR : %v\n",err.Error())
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	dk.Client.ContainerCreate(
		dk.Context,
		&container.Config{
			User : "root",
			ExposedPorts: nat.PortSet{
				"80" : struct{}{},
			},
		},
		&container.HostConfig{},
		&network.NetworkingConfig{},
		&v1.Platform{},
		"test",
	)
	
	return
}
*/