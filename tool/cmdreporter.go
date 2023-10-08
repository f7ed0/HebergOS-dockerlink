package tool

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"runtime"

	"github.com/f7ed0/HebergOS-dockerlink/docker"

	"github.com/docker/docker/api/types"
)

func CmdReporter(resp http.ResponseWriter,req *http.Request,container_id string,command []string,haspath bool,path string) {


	exarg := types.ExecConfig{
		User : "root",
		AttachStderr: true,
		AttachStdout: true,
		Cmd: command,
	}
	if haspath {
		exarg.WorkingDir = path
	}

	dk,err := docker.NewDockerHandler()
	if err != nil {
		log.Default().Println(err.Error())
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}
	idresp,err := dk.Client.ContainerExecCreate(dk.Context,container_id,exarg)
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
		resp.Write([]byte(res[8:]))
		return
	}

	fmt.Println(res)
	resp.Write([]byte(res[8:]))
	return
}