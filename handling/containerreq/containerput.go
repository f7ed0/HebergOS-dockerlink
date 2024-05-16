package containerreq

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"
	"os/exec"
	"strconv"

	"github.com/f7ed0/HebergOS-dockerlink/consts"
	"github.com/f7ed0/HebergOS-dockerlink/docker"
	"github.com/f7ed0/HebergOS-dockerlink/logger"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
)

func ContainerPut(resp http.ResponseWriter, req *http.Request) {

	PORT_AREA_SIZE, err := strconv.Atoi(os.Getenv("port_area_size"))
	if err != nil {
		panic("port_area_size not present or malformed")
	}

	resp.Header().Set("Content-Type", "text/plain")

	var p map[string]any

	jrd := json.NewDecoder(req.Body)

	err = jrd.Decode(&p)
	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Println(p)

	// Get ports range
	name, ok := p["name"].(string)
	if !ok {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}
	ports, ok := p["host_port_root"].(float64)
	if !ok {
		resp.WriteHeader(http.StatusBadRequest)
		resp.Write([]byte("missing host_port_root"))
		return
	}
	ports_int := int(ports)

	// Get memory limit
	mem, ok := p["memory"].(float64)
	if !ok {
		resp.WriteHeader(http.StatusBadRequest)
		resp.Write([]byte("missing memory"))
		return
	}

	// GetCPUlimit
	cpu, ok := p["cpulimit"].(float64)
	if !ok {
		resp.WriteHeader(http.StatusBadRequest)
		resp.Write([]byte("missing cpulimit"))
		return
	}

	// Get Image id
	img, ok := p["image"].(string)
	if !ok {
		img = os.Getenv("base_image")
	}

	command_concat := "sh /var/www/starter.sh"

	// Get stating command
	cmds, ok := p["commands"].([]string)
	if ok {
		for _, cmd := range cmds {
			command_concat += " && " + cmd
		}
	}

	command_concat += " && /usr/sbin/sshd -D"

	// Get ports to forward
	prts := nat.PortSet{
		"22/tcp":  {},
		"80/tcp":  {},
		"443/tcp": {},
	}
	np, ok := p["ports"].([]any)
	logger.Default.Log("INFO", "%v %v %v", p["ports"], np, ok)
	if ok {
		var pt nat.Port
		for _, port := range np {
			ui := true
			for k := range prts {
				pt = nat.Port(port.(string))
				if pt.Int() >= PORT_AREA_SIZE {
					resp.WriteHeader(http.StatusBadRequest)
					resp.Write([]byte("port area size is " + os.Getenv("port_area_size")))
					return
				}
				if k.Int() == pt.Int() {
					ui = false
				}
			}
			logger.Default.Log("INFO", "%v %v", ui, pt)
			if !ui {
				continue
			}
			prts[pt] = struct{}{}
		}
	}
	prtmap := nat.PortMap{}
	for k := range prts {
		prtmap[k] = []nat.PortBinding{
			{
				HostIP:   "0.0.0.0",
				HostPort: strconv.Itoa(k.Int() + ports_int),
			},
		}
	}

	// Creating the container
	j := json.NewEncoder(resp)

	file, err := os.Create(os.Getenv("nginxconfdir") + "/sites-available/" + name + ".conf")
	if err != nil {
		logger.Default.Log("ERR", err.Error())
		resp.WriteHeader(http.StatusPreconditionFailed)
		resp.Write([]byte(err.Error()))
		return
	}

	_, err = fmt.Fprintf(file, consts.NGINX_TEMPLATE, name, ports_int+80, name, ports_int)
	if err != nil {
		logger.Default.Log("ERR", err.Error())
		resp.WriteHeader(http.StatusPreconditionFailed)
		resp.Write([]byte(err.Error()))
		return
	}

	file.Close()

	cmd := exec.Command("ln", "-s", os.Getenv("nginxconfdir")+"/sites-available/"+name+".conf", os.Getenv("nginxconfdir")+"/sites-enabled/"+name+".conf")
	if err := cmd.Run(); err != nil {
		if !(cmd.ProcessState.ExitCode() == 1) {
			res, err2 := cmd.Output()
			if err2 != nil {
				logger.Default.Log("ERR", err.Error()+" @ ln -s")
			} else {
				logger.Default.Log("ERR", err.Error()+" @ ln -s\n"+string(res))
			}
			resp.WriteHeader(http.StatusPreconditionFailed)

			resp.Write([]byte(err.Error()))
			return
		}
	}

	cmd = exec.Command("/usr/sbin/nginx", "-t")
	if err := cmd.Run(); err != nil {
		res, err2 := cmd.Output()
		if err2 != nil {
			logger.Default.Log("ERR", err.Error()+" @ nginx -t")
		}
		logger.Default.Log("ERR", err.Error()+" @ nginx -t\n"+string(res))
		resp.WriteHeader(http.StatusPreconditionFailed)
		resp.Write([]byte(err.Error() + " @ nginx -t"))
		return
	}
	cmd = exec.Command("systemctl", "reload", "nginx")
	if err := cmd.Run(); err != nil {
		res, err2 := cmd.Output()
		if err2 != nil {
			logger.Default.Log("ERR", err.Error()+" @ systemctl reload nginx")
		} else {
			logger.Default.Log("ERR", err.Error()+" @ systemctl reload nginx\n"+string(res))
		}
		resp.WriteHeader(http.StatusPreconditionFailed)
		resp.Write([]byte(err.Error() + " @ systemctl reload nginx"))
		return
	}
	cmd = exec.Command("certbot", "--nginx", "-d", "ssh."+name+".insash.org", "-n")
	if err := cmd.Run(); err != nil {
		res, err2 := cmd.Output()
		if err2 != nil {
			logger.Default.Log("ERR", err.Error()+" @ certbot --nginx -d -n")
		} else {
			logger.Default.Log("ERR", err.Error()+" @ certbot --nginx -d -n\n"+string(res))
		}
		resp.WriteHeader(http.StatusPreconditionFailed)
		resp.Write([]byte(err.Error() + " @ certbot --nginx -d"))
		return
	}

	dk, err := docker.NewDockerHandler()
	if err != nil {
		logger.Default.Log("ERR", err.Error())
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer dk.Client.Close()

	// Creating container
	client, err := dk.Client.ContainerCreate(
		dk.Context,
		&container.Config{
			User:         "root",
			Image:        img,
			WorkingDir:   "/var/www/html",
			Cmd:          []string{"sh", "-c", command_concat},
			ExposedPorts: prts,
			Labels:       map[string]string{"ports": strconv.Itoa(ports_int), "dockerlink": consts.DOCKERLINK_VERSION},
			Volumes: map[string]struct{}{
				"/var/www": {},
			},
		},
		&container.HostConfig{
			PortBindings: prtmap,
			Resources: container.Resources{
				Memory:    int64(mem * math.Pow(2, 30)),
				CPUQuota:  int64(cpu * 10000),
				CPUPeriod: 10000,
			},
			Links: []string{os.Getenv("mysql_docker_name"), os.Getenv("postgre_docker_name")},
		},
		&network.NetworkingConfig{},
		&v1.Platform{
			OS: "linux",
		},
		name,
	)

	if err != nil {
		logger.Default.Log("ERR", err.Error())
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(err.Error()))
		return
	}

	resp.Header().Set("Content-Type", "application/json")

	j.Encode(client)

	return
}
