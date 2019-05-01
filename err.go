package main

import "errors"

var (
	errCommandNotFound       = errors.New("command cannot not be found")
	errCommandNotRegularFile = errors.New("command is not a regular file")
	errCommandHiddenFile     = errors.New("command cannot be a hidden file")
	errCommandNotExecutable  = errors.New("command is not executable")
	errZeroIterations        = errors.New("iterations cannot be zero")
	errGobinNotFound         = errors.New("GOBIN or GOPATH environment variable not set")
	errBasePathNotRecovered  = errors.New("cannot recover the base path")
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
