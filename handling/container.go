package handling

import (
	"log"
	"net/http"

	"github.com/f7ed0/HebergOS-dockerlink/handling/containerreq"
)

const BASE_IMAGE = "41287f5341c2713f9d444f3d55fec01bae3ffd9f5302f65dc9747caf6aba32fc"

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