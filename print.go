package main

import (
	"fmt"
)

func print(s string) {
	fmt.Println("**\t" + s)
}

func printf(pattern string, s ...interface{}) {
	print(fmt.Sprintf(pattern, s...))
}
