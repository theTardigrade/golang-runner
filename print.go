package main

import (
	"fmt"
	"strings"
)

var (
	fourSpaces = strings.Repeat(" ", 4)
)

func print(s string) {
	fmt.Printf("%s[%d]%s%s\n", fourSpaces, unixMilli(), fourSpaces, s)
}

func printf(pattern string, s ...interface{}) {
	print(fmt.Sprintf(pattern, s...))
}
