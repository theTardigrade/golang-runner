package main

import (
	"fmt"
	"runtime"
	"strings"
)

var (
	fourSpaces = strings.Repeat(" ", 4)
	newline    string
)

func init() {
	if runtime.GOOS == "windows" {
		newline = "\r\n"
	} else {
		newline = "\n"
	}
}

func print(s ...string) {
	fmt.Printf("%s[%d]%s%s\n", fourSpaces, unixMilli(), fourSpaces, strings.Join(s, " "))
}

func printf(pattern string, s ...interface{}) {
	print(fmt.Sprintf(pattern, s...))
}
