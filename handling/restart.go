package handling

import (
	"net/http"

	"github.com/docker/docker/api/types/container"
	"github.com/f7ed0/HebergOS-dockerlink/docker"
	"github.com/f7ed0/HebergOS-dockerlink/logger"
	"github.com/f7ed0/HebergOS-dockerlink/tool"
)

func RestartDocker(resp http.ResponseWriter,req *http.Request) {
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
		logger.Default.Log("ERR",err.Error())
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	info,err := dk.Client.ContainerInspect(dk.Context,id[0])

	if err != nil {
		resp.WriteHeader(http.StatusPreconditionFailed)
		resp.Write([]byte(err.Error()))
		return
	}

	if !info.State.Running {
		resp.WriteHeader(http.StatusNotAcceptable)
		resp.Write([]byte("Not running"))
		return
	}

	dkl_version,ok := info.Config.Labels["dockerlink"]
	if !ok {
		resp.WriteHeader(http.StatusPreconditionFailed)
		resp.Write([]byte("this docker has not been created with dockerlink"))
		return
	}

	if !tool.VersionCheck(dkl_version,"v0.0",tool.VersionSup) {
		resp.WriteHeader(http.StatusPreconditionFailed)
		resp.Write([]byte("this docker has not been created with dockerlink or is too old"))
		return
	}

	t := 0
	err = dk.Client.ContainerRestart(dk.Context,id[0],container.StopOptions{
		Timeout: &t,
	})
	if err != nil {
		logger.Default.Log("ERR",err.Error())
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Header().Set("Content-Type", "text/plain")
		resp.Write([]byte(err.Error()))
		return
	}
		
	resp.WriteHeader(http.StatusNoContent)
	return
}