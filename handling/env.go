package handling

import (
	"encoding/json"
	"net/http"
	"os"
)

func Env(resp http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		j := json.NewEncoder(resp)
		resp.Header().Set("Content-Type", "application/json")

		j.Encode(map[string]string{
			"docker_sock" : os.Getenv("docker_sock"),
			"port_area_size" : os.Getenv("port_area_size"),
			"base_image" : os.Getenv("base_image"),
		})
		return
	}
	resp.WriteHeader(http.StatusMethodNotAllowed)
}