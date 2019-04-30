package main

import "errors"

var (
	errCommandNotFound = errors.New("command not found")
	errZeroIterations  = errors.New("zero iterations")
	errGobinNotFound   = errors.New("GOBIN or GOPATH environment variable not set")
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func judgeErr(err error) string {
	if err == nil {
		return "SUCCESS"
	}

	return "FAILURE"
}
