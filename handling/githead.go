package handling

import (
	"herbergOS/tool"
	"net/http"
)

func GitHead(resp http.ResponseWriter,req *http.Request) {
	if(req.Method == http.MethodGet) {
		gitHeadGet(resp,req)
		return
	}
	resp.WriteHeader(http.StatusMethodNotAllowed)
}

func gitHeadGet(resp http.ResponseWriter, req *http.Request) {
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

	tool.CmdReporter(resp,req,id[0],[]string{"git","log","-1"},ok,path[0])
}