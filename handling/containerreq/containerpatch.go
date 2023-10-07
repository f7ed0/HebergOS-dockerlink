package containerreq

import (
	"encoding/json"
	"herbergOS/docker"
	"log"
	"math"
	"net/http"

	"github.com/docker/docker/api/types/container"
)

func ContainerPatch(resp http.ResponseWriter, req *http.Request) {
	qmap := req.URL.Query()

	id,ok := qmap["id"]
	if !ok {
		resp.WriteHeader(http.StatusBadRequest)
		resp.Write([]byte("id is missing"))
		return
	}

	dk,err := docker.NewDockerHandler()
	if err != nil {
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

	res := info.HostConfig.Resources


	var p map[string]any

	jrd := json.NewDecoder(req.Body)

	err = jrd.Decode(&p)
	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	cpulimit,ok := p["cpulimit"].(float64)
	if ok {
		res.CPUQuota = int64(cpulimit*10000)
	}

	memlimit,ok := p["memory"].(float64)
	if ok {
		res.Memory = int64(memlimit*math.Pow(2,30))
		log.Default().Println(int64(memlimit*math.Pow(2,30)))
	}

	js := json.NewEncoder(log.Default().Writer())

	js.Encode(res)

	update,err := dk.Client.ContainerUpdate(dk.Context,id[0],container.UpdateConfig{
		Resources: res,
	})
	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		resp.Write([]byte(err.Error()))
		return
	}
	
	j := json.NewEncoder(resp)
	err = j.Encode(update)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(err.Error()))
		return
	}
}