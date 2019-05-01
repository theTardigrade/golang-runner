package main

import (
	"os"
	"os/signal"
	"strings"
)

func init() {
	go watchSignals()
}

// runs in own goroutine
func watchSignals() {
	ch := make(chan os.Signal, 1)

	signal.Notify(ch, os.Interrupt, os.Kill)

	if s := <-ch; *flagVerbose {
		printf("%s SIGNAL RECEIVED", strings.ToUpper(s.String()))
	}

	exit()
}
