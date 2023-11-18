package containerreq

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/f7ed0/HebergOS-dockerlink/docker"
	"github.com/f7ed0/HebergOS-dockerlink/logger"

	"github.com/docker/docker/api/types"
)

func ContainerGet(resp http.ResponseWriter,req *http.Request) {
	qmap := req.URL.Query()
	j := json.NewEncoder(resp)

	ids,ok := qmap["id"]

	dk,err := docker.NewDockerHandler()
	if err != nil {
		logger.Default.Log("ERR",err.Error())
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
			logger.Default.Log("ERR",err.Error())
			resp.WriteHeader(http.StatusInternalServerError)
			return
		}
	
		ids = []string{}
	
		for _,container := range containers {
			ids = append(ids, container.ID)
		}
	}

	result := map[string]map[string]any{"success" : {}, "error": {}}
	for _,id := range ids {
		json,err := dk.Client.ContainerInspect(dk.Context,id)
		if err != nil {
			logger.Default.Log("ERR",err.Error())
			result["error"][id] = err.Error()
		} else {
			if  has_run && json.State.Running != running {
				continue
			}
	
			result["success"][id] = map[string]any{"name":json.Name,"state":json.State.Status}
			
			vers,ok := json.Config.Labels["dockerlink"]
			if ok {
				result["success"][id].(map[string]any)["dockerlink"] = vers
			} else {
				result["success"][id].(map[string]any)["dockerlink"] = "not_dockerlink"
			}
		
			result["success"][id].(map[string]any)["host_port_root"] = json.Config.Labels["ports"]
			
			result["success"][id].(map[string]any)["ports"] = json.Config.ExposedPorts
			
			if json.State.Running {
				t,err := time.Parse(time.RFC3339Nano,json.State.StartedAt)
				if err != nil {
					logger.Default.Log("ERR",err.Error())
					resp.WriteHeader(http.StatusInternalServerError)
					return
				}
		
				result["success"][id].(map[string]any)["started_at"] = t.Unix()
			} else {
				result["success"][id].(map[string]any)["exit_code"] = json.State.ExitCode
			}
		}
	}

	resp.Header().Set("Content-Type", "application/json")
	j.Encode(result)

	return
}