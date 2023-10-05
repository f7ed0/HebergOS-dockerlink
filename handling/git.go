package handling

import (
	"herbergOS/tool"
	"net/http"
)

func Git(resp http.ResponseWriter,req *http.Request) {
	if(req.Method == http.MethodGet) {
		gitGet(resp,req)
		return
	}
	if(req.Method == http.MethodPut) {
		gitPut(resp,req)
		return
	}
	resp.WriteHeader(http.StatusMethodNotAllowed)
	return
}

func gitGet(resp http.ResponseWriter, req *http.Request) {
	qmap := req.URL.Query()

	id,ok := qmap["id"]
	if(!ok) {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	path,ok := qmap["path"]

	if !ok {
		path = []string{""}
	}

	tool.CmdReporter(resp,req,id[0],[]string{"git","pull"},ok,path[0])

}

func gitPut(resp http.ResponseWriter, req *http.Request) {
	qmap := req.URL.Query()

	id,ok := qmap["id"]
	if(!ok) {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	path,ok := qmap["path"]

	if !ok {
		path = []string{""}
	}

	tool.CmdReporter(resp,req,id[0],[]string{"git","init"},ok,path[0])
}