package tool

import (
	"strconv"
	"strings"
)

type Versionchecker func([]int,[]int) bool

func VersionSup(current_version,required_version []int) bool {
	for i := 0 ; i < len(current_version) && i < len(required_version) ; i ++ {
		if(current_version[i] > required_version[i]) {
			return true
		}
		if(current_version[i] < required_version[i]) {
			return false
		}
	}
	if len(current_version) >= len(required_version) {
		return true
	}
	return false
}

func VersionCheck(current_version,required_version string, check Versionchecker) bool {
	if !strings.HasPrefix(current_version,"v") || !strings.HasPrefix(required_version,"v") {
		return false
	}
	curr_split := strings.Split(current_version[1:],".")
	req_split := strings.Split(required_version[1:],".")

	curr_split_int := []int{}
	req_split_int := []int{}

	for _,subversion := range curr_split {
		subint,err := strconv.Atoi(subversion)
		if err != nil {
			return false
		}
		curr_split_int = append(curr_split_int,subint)
	}
	for _,subversion := range req_split {
		subint,err := strconv.Atoi(subversion)
		if err != nil {
			return false
		}
		req_split_int = append(req_split_int,subint)
	}
	return check(curr_split_int,req_split_int)
}