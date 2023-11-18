package handling

import (
	"net/http"

	"github.com/f7ed0/HebergOS-dockerlink/handling/containerreq"
	"github.com/f7ed0/HebergOS-dockerlink/logger"
)

func Container(resp http.ResponseWriter,req *http.Request) {
	logger.Default.Log("REQ","%s => %s",req.Method,req.URL.RequestURI())
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