package handling

import (
	"fmt"
	"herbergOS/docker"
	"net/http"
	"strconv"
)


func Stats(resp http.ResponseWriter,req *http.Request) {
	if(req.Method != "GET") {
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
	}

	intsince,err := strconv.ParseInt(since[0],10,64)
	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)
	}

	resp.Header().Set("Content-Type", "application/json")

	fmt.Fprint(resp,docker.Sh.Export(id[0],intsince))
	
	return
}