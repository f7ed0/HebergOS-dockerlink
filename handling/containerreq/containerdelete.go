package containerreq

import (
	"log"
	"net/http"

	"github.com/f7ed0/HebergOS-dockerlink/docker"

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

	docker.Sh.Wipe(id[0])

	resp.WriteHeader(http.StatusNoContent)
}