package handling

import (
	"log"
	"net/http"

	"github.com/f7ed0/HebergOS-dockerlink/handling/containerreq"
)

func Container(resp http.ResponseWriter,req *http.Request) {
	log.Default().Println(req.Method)
	switch(req.Method) {
	case http.MethodGet:
		containerreq.ContainerGet(resp,req)
		return
	case http.MethodPut:
		containerreq.ContainerPut(resp,req)
		return
	case http.MethodDelete:
		containerreq.ContainerDelete(resp,req)
		return
	case http.MethodPatch:
		containerreq.ContainerPatch(resp,req)
		return
	}
	resp.WriteHeader(http.StatusMethodNotAllowed)
	return
}