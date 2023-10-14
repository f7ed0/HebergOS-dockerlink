package handling

import (
	"net/http"

	"github.com/f7ed0/HebergOS-dockerlink/tool"
)

func GitBranch(resp http.ResponseWriter,req *http.Request) {
	if(req.Method == http.MethodGet) {
		gitBranchGet(resp,req)
		return
	}
	if(req.Method == "POST"){
		gitBranchPost(resp,req)
		return
	}
	resp.WriteHeader(http.StatusMethodNotAllowed)
}

func gitBranchGet(resp http.ResponseWriter,req *http.Request) {
	qmap := req.URL.Query()

	id,ok := qmap["id"]
	if !ok {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	path,ok := qmap["path"]

	if !ok {
		path = []string{""}
	}

	tool.CmdReporter(resp,req,id[0],[]string{"git","rev-parse","--abbrev-ref","HEAD"},ok,path[0],"admin")
}

func gitBranchPost(resp http.ResponseWriter, req *http.Request) {
	qmap := req.URL.Query()

	id,ok := qmap["id"]
	if !ok {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}
	branch,ok := qmap["branch"]
	if !ok {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	path,ok := qmap["path"]

	if !ok {
		path = []string{""}
	}

	tool.CmdReporter(resp,req,id[0],[]string{"git","checkout",branch[0]},ok,path[0],"admin")
}