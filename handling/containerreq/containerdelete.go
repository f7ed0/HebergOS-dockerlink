package containerreq

import (
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/f7ed0/HebergOS-dockerlink/docker"
	"github.com/f7ed0/HebergOS-dockerlink/tool"

	"github.com/docker/docker/api/types"
)

func ContainerDelete(resp http.ResponseWriter,req *http.Request) {
	qmap := req.URL.Query()

	resp.Header().Set("Content-Type", "text/plain")

	id,ok := qmap["id"]
	if !ok {
		resp.WriteHeader(http.StatusBadRequest)
		resp.Write([]byte("missing id"))
		return
	}

	var delete_volume = true
	dv,ok := qmap["deletevolume"]
	if ok {
		delete_volume = dv[0] == "true"
	}

	dk,err := docker.NewDockerHandler()
	if err != nil {
		log.Default().Println(err.Error())
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(err.Error()))
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

	

	err = dk.Client.ContainerRemove(dk.Context,id[0],types.ContainerRemoveOptions{
		RemoveVolumes: delete_volume,
		RemoveLinks: false,
		Force: true,
	})
	if err != nil {
		log.Default().Println(err.Error())
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(err.Error()))
		return
	}

	cmd := exec.Command("screen","-X","-S",info.Name+"wettyssh","quit")
	cmd.Dir = os.Getenv("wettydir")
	if err := cmd.Run(); err != nil {
		log.Default().Println(err.Error())
		resp.WriteHeader(http.StatusAccepted)
		resp.Write([]byte(err.Error()))
		return
	}

	resp.WriteHeader(http.StatusNoContent)
}