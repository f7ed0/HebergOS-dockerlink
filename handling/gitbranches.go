package handling

import (
	"bytes"
	"encoding/json"
	"herbergOS/docker"
	"io"
	"log"
	"net/http"
	"runtime"
	"strings"

	"github.com/docker/docker/api/types"
)

func GitBranches(resp http.ResponseWriter, req *http.Request) {
	if(req.Method == http.MethodGet) {
		gitBranchesGet(resp,req)
		return
	}
	resp.WriteHeader(http.StatusMethodNotAllowed)
}

func gitBranchesGet(resp http.ResponseWriter,req *http.Request) {
	qmap := req.URL.Query()

	id,ok := qmap["id"]
	if(!ok) {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	path,ok := qmap["path"]

	// Preparing the execution of the command
	exarg := types.ExecConfig{
		User : "root",
		AttachStderr: true,
		AttachStdout: true,
		Cmd: []string{"git","branch","-a"},
	}
	if ok {
		exarg.WorkingDir = path[0]
	}

	// Create the execution
	jsonw := json.NewEncoder(resp)
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
	
	l := []string{}

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
	if retcode.ExitCode != 0 {
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
			res += string(bytes.ToValidUTF8(str,[]byte("")))
		}
		resp.WriteHeader(http.StatusPreconditionFailed)
		resp.Header().Set("Content-Type", "text/plain")
		resp.Write([]byte(res[7:]))
		return
	}
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
		x := strings.Split(string(str)," ")
		l = append(l, x[len(x)-1])
	}
	resp.Header().Set("Content-Type", "application/json")
	jsonw.Encode(l)
}