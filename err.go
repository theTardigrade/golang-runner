package main

import "errors"

var (
	errCommandNotFound = errors.New("command not found")
	errZeroIterations  = errors.New("zero iterations")
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
