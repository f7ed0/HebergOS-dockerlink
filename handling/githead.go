package handling

import (
	"net/http"

	"github.com/f7ed0/HebergOS-dockerlink/docker"
	"github.com/f7ed0/HebergOS-dockerlink/tool"
)

func GitHead(resp http.ResponseWriter,req *http.Request) {
	if(req.Method == http.MethodGet) {
		gitHeadGet(resp,req)
		return
	}
	resp.WriteHeader(http.StatusMethodNotAllowed)
}

func gitHeadGet(resp http.ResponseWriter, req *http.Request) {
	qmap := req.URL.Query()

	id,ok := qmap["id"]
	if(!ok) {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	dk,err := docker.NewDockerHandler()

	info,err := dk.Client.ContainerInspect(dk.Context,id[0])

	if err != nil {
		resp.WriteHeader(http.StatusPreconditionFailed)
		resp.Write([]byte(err.Error()))
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

	path,ok := qmap["path"]

	if !ok {
		path = []string{""}
	}

	tool.CmdReporter(resp,req,id[0],[]string{"git","log","-1"},ok,path[0],"admin")
}