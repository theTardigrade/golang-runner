package main

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
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
	_, path, _, ok := runtime.Caller(0)
	if !ok {
		panic(errBasePathNotRecovered)
	}
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

func isPathHidden(path string) bool {
	name := filepath.Base(path)

	if strings.HasPrefix(name, pathHiddenNamePrefix) || strings.HasSuffix(name, pathHiddenNameSuffix) {
		return true
	}

	return false
}
