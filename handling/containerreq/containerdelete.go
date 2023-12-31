package containerreq

import (
	"net/http"
	"os"
	"os/exec"

	"github.com/f7ed0/HebergOS-dockerlink/docker"
	"github.com/f7ed0/HebergOS-dockerlink/logger"
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
		logger.Default.Log("ERR",err.Error())
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(err.Error()))
		return
	}

	info,err := dk.Client.ContainerInspect(dk.Context,id[0])
	if err != nil {
		logger.Default.Log("ERR",err.Error())
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

	if info.State.Running {
		resp.WriteHeader(http.StatusPreconditionFailed)
		resp.Write([]byte("docker needs to be stopped before deleting it"))
		return
	}

	cmd := exec.Command("rm",os.Getenv("nginxconfdir")+"/sites-available/"+info.Name+".conf",os.Getenv("nginxconfdir")+"/sites-enabled/"+info.Name+".conf")
	if err := cmd.Run(); err != nil {
		if !(cmd.ProcessState.ExitCode() == 1) {
			res,err2 := cmd.Output()
			if err2 != nil {
				logger.Default.Log("ERR",err.Error()+" @ rm")
			} else {
				logger.Default.Log("ERR",err.Error()+" @ rm\n"+string(res))
			}
			resp.WriteHeader(http.StatusPreconditionFailed)
			
			resp.Write([]byte(err.Error()+" @ rm"))
			return
		}
	}

	cmd = exec.Command("/usr/sbin/nginx","-t")
	if err := cmd.Run(); err != nil {
		res,err2 := cmd.Output()
		if err2 != nil {
			logger.Default.Log("ERR",err.Error()+" @ nginx -t")
		} else {
			logger.Default.Log("ERR",err.Error()+" @ nginx -t\n"+string(res))
		}
		resp.WriteHeader(http.StatusPreconditionFailed)
		resp.Write([]byte(err.Error()+" @ nginx -t"))
		return
	}
	cmd = exec.Command("systemctl","reload","nginx")
	if err := cmd.Run(); err != nil {
		res,err2 := cmd.Output()
		if err2 != nil {
			logger.Default.Log("ERR",err.Error()+" @ systemctl reload nginx")
		} else {
			logger.Default.Log("ERR",err.Error()+" @ systemctl reload nginx\n"+string(res))
		}
		resp.WriteHeader(http.StatusPreconditionFailed)
		resp.Write([]byte(err.Error()+" @ systemctl reload nginx"))
		return
	}

	err = dk.Client.ContainerRemove(dk.Context,id[0],types.ContainerRemoveOptions{
		RemoveVolumes: delete_volume,
		RemoveLinks: false,
		Force: true,
	})
	if err != nil {
		logger.Default.Log("ERR",err.Error())
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(err.Error()))
		return
	}

	docker.Sh.Wipe(id[0])

	resp.WriteHeader(http.StatusNoContent)
}