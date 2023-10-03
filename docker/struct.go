package docker

import (
	"context"
	"os"

	"github.com/docker/docker/client"
)

type DockerHandler struct {
	Context context.Context
	Client *client.Client
}

func NewDockerHandler(ops ...client.Opt) (DockerHandler,error) {
	ctx := context.Background()
	cli,err := client.NewClientWithOpts(append(ops, client.WithHost(os.Getenv("docker_sock")))...)
	return DockerHandler{Context: ctx,Client:cli},err
}