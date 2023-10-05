package handling

import (
	"bytes"
	"herbergOS/docker"
	"io"
	"log"
	"net/http"
	"runtime"

	"github.com/docker/docker/api/types"
)

func Git(resp http.ResponseWriter,req *http.Request) {
	if(req.Method == "GET") {
		gitGet(resp,req)
		return
	}
	resp.WriteHeader(http.StatusMethodNotAllowed)
}

func gitGet(resp http.ResponseWriter, req *http.Request) {
	qmap := req.URL.Query()

	id,ok := qmap["id"]
	if(!ok) {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	path,ok := qmap["path"]

	exarg := types.ExecConfig{
		User : "root",
		AttachStderr: true,
		AttachStdout: true,
		Cmd: []string{"git","fetch","--all"},
	}
	if ok {
		exarg.WorkingDir = path[0]
	}

	dk,err := docker.NewDockerHandler()
	if err != nil {
		log.Default().Println(err.Error())
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}
	idresp,err := dk.Client.ContainerExecCreate(dk.Context,id[0],exarg)
	if err != nil {
		log.Default().Println(err.Error())
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Execute and attach stdout and stderr to a reader
	response,err := dk.Client.ContainerExecAttach(dk.Context,idresp.ID,types.ExecStartCheck{

	})
	if err != nil {
		log.Default().Println(err.Error())
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer response.Close()

	// inspect the exec to get the return code
	retcode,err := dk.Client.ContainerExecInspect(dk.Context,idresp.ID)
	if err != nil {
		log.Default().Println(err.Error())
		resp.WriteHeader(http.StatusInternalServerError)
	}
	for retcode.Running {
		retcode,err = dk.Client.ContainerExecInspect(dk.Context,idresp.ID)
		if err != nil {
			log.Default().Println(err.Error())
			resp.WriteHeader(http.StatusInternalServerError)
			return
		}
		runtime.Gosched()
	}

	// sending data
	res := ""
	for true {
		str,_,err := response.Reader.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Default().Println(err.Error())
			resp.WriteHeader(http.StatusInternalServerError)
			return
		}
		res += string(bytes.ToValidUTF8(str,[]byte("")))+"\n"
	}

	// sending data
	resp.Header().Set("Content-Type", "text/plain")
	if retcode.ExitCode != 0 {
		log.Default().Println("err")
		resp.WriteHeader(http.StatusPreconditionFailed)
		resp.Write([]byte(res[7:]))
		return
	}
	
	resp.Write([]byte(res[:7]))
	return
}