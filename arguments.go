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
		splitArguments := strings.Split(*internalFlag.Arguments, " ")

		if l := len(splitArguments); l > 0 {
			arguments = make([]string, 0, l)

			for i := 0; i < l; i++ {
				if a := splitArguments[i]; a != "" {
					arguments = append(arguments, a)
				}
			}
		}
	}
}
