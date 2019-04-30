package main

import (
	"os"
	"os/signal"
)

func init() {
	go func() {
		ch := make(chan os.Signal, 1)

		signal.Notify(ch, os.Interrupt, os.Kill)

		<-ch
		exit()
	}()
}
