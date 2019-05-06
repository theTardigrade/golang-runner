package main

import (
	"strings"

	internalFlag "github.com/theTardigrade/runner/internal/flag"
)

var (
	arguments []string
)

func init() {
	if *internalFlag.Arguments != "" {
		for _, a := range strings.Split(*internalFlag.Arguments, " ") {
			if a != "" {
				arguments = append(arguments, a)
			}
		}
	}
}
