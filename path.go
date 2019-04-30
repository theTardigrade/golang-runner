package main

import (
	"os"
	"path/filepath"
	"runtime"
)

const (
	pathWindowsNameSuffix = ".exe"
	pathHiddenNameSuffix  = "~"
	pathHiddenNamePrefix  = "."
)

var (
	basePath  string
	gobinPath string
)

func init() {
	_, path, _, _ := runtime.Caller(0)
	basePath = filepath.Dir(path)

	gobinPath = gobin()
}

func gobin() (value string) {
	value, found := os.LookupEnv("GOBIN")

	if !found {
		value, found = os.LookupEnv("GOPATH")
		if !found {
			panic(errGobinNotFound)
		}
		value = filepath.Join(value, "bin")
	}

	return
}
