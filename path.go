package main

import (
	"path/filepath"

	"github.com/theTardigrade/envStore"
)

const (
	pathWindowsNameSuffix = ".exe"
	pathHiddenNameSuffix  = "~"
	pathHiddenNamePrefix  = "."
)

func gobin() (value string) {
	env, err := envStore.New(envStore.Config{
		FromSystem: true,
	})
	checkErr(err)

	value, err = env.Get("GOBIN")
	if err != envStore.ErrKeyNotFound {
		checkErr(err)
	}

	value, err = env.Get("GOPATH")
	checkErr(err)
	value = filepath.Join(value, "bin")

	return
}
