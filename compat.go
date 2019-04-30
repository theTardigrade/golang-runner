package main

import (
	"runtime"
)

var (
	isWindows bool
	newline   string
)

func init() {
	isWindows = (runtime.GOOS == "windows")

	if isWindows {
		newline = "\r\n"
	} else {
		newline = "\n"
	}
}
