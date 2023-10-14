package handling

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/f7ed0/HebergOS-dockerlink/docker"
)


func Stats(resp http.ResponseWriter,req *http.Request) {
	if(req.Method != http.MethodGet) {
		resp.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	Qmap := req.URL.Query()

	id,ok := Qmap["id"]
	if(!ok) {
		resp.WriteHeader(http.StatusBadRequest)
	}

	since,ok := Qmap["since"]
	if(!ok) {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	scale,ok := Qmap["scale"]
	if(!ok) {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}
	only := docker.INSTANT_ONLY
	if (scale[0] == "day"){
		only = docker.DAY_ONLY
	} else if (scale[0] == "week") {
		only = docker.WEEK_ONLY
	}

	intsince,err := strconv.ParseInt(since[0],10,64)
	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	resp.Header().Set("Content-Type", "application/json")

	fmt.Fprint(resp,docker.Sh.Export(id[0],intsince,only))

	return
}