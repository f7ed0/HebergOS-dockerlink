package handling

import (
	"net/http"
	"os"
	"os/exec"
	"strconv"

	"github.com/f7ed0/HebergOS-dockerlink/docker"
	"github.com/f7ed0/HebergOS-dockerlink/logger"
	"github.com/f7ed0/HebergOS-dockerlink/tool"

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

	if info.State.Running {
		resp.WriteHeader(http.StatusNotAcceptable)
		resp.Write([]byte("Already running"))
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

	err = dk.Client.ContainerStart(dk.Context,id[0],types.ContainerStartOptions{})
	if err != nil {
		logger.Default.Log("ERR",err.Error())
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Header().Set("Content-Type", "text/plain")
		resp.Write([]byte(err.Error()))
		return
	}

	port,ok := info.Config.Labels["ports"]

	if !ok {
		resp.WriteHeader(http.StatusAccepted)
		return
	}

	intport,err := strconv.Atoi(port)

	if !ok {
		resp.WriteHeader(http.StatusAccepted)
		return
	}

	cmd := exec.Command("screen","-dmS",info.Name+"wettyssh","/usr/bin/node",".","--ssh-host=localhost","--ssh-port="+strconv.Itoa(intport+22),"--port="+strconv.Itoa(intport),"--force-ssh","--allow-iframe","--bypass-helmet")
	cmd.Dir = os.Getenv("wettydir")
	if err := cmd.Run(); err != nil {
		logger.Default.Log("ERR",err.Error())
		resp.WriteHeader(http.StatusAccepted)
		resp.Write([]byte(err.Error()))
		return
	}
		
	resp.WriteHeader(http.StatusNoContent)
	return
}