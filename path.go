package main

import (
	"os"
	"path/filepath"
)

const (
	pathWindowsNameSuffix = ".exe"
	pathHiddenNameSuffix  = "~"
	pathHiddenNamePrefix  = "."
)

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
