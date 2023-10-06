package handling

import (
	"herbergOS/docker"
	"log"
	"net/http"
	"runtime"

	"github.com/docker/docker/api/types"
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

	r,err := dk.Client.ContainerExecCreate(dk.Context,id[0],types.ExecConfig{
		User: "root",
		Cmd: []string{"pkill","screen"},
	})
	if err != nil {
		panic(err)
	}
	err = dk.Client.ContainerExecStart(dk.Context,r.ID,types.ExecStartCheck{
	})
	if err != nil {
		panic(err)
	}
	info,err := dk.Client.ContainerExecInspect(dk.Context,r.ID)
	if err != nil {
		panic(err)
	}

	for info.Running {
		runtime.Gosched()
	}

	t := 0
	log.Default().Println(info.ExitCode)
	err = dk.Client.ContainerStop(dk.Context,id[0],container.StopOptions{
		Timeout: &t,
	})
	
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