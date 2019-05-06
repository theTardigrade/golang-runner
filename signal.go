package main

import (
	"os"
	"os/signal"
	"strings"

	internalFlag "github.com/theTardigrade/runner/internal/flag"
	internalFmt "github.com/theTardigrade/runner/internal/fmt"
)

func init() {
	go watchSignals()
}

// runs in own goroutine
func watchSignals() {
	ch := make(chan os.Signal, 1)

	signal.Notify(ch, os.Interrupt, os.Kill)

	if s := <-ch; *internalFlag.Verbose {
		internalFmt.Printf("%s SIGNAL RECEIVED", strings.ToUpper(s.String()))
	}

	exit()
}
