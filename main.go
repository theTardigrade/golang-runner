package main

import (
	"context"
	"flag"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	internalErrors "github.com/theTardigrade/runner/internal/errors"
	internalFlag "github.com/theTardigrade/runner/internal/flag"
	internalFmt "github.com/theTardigrade/runner/internal/fmt"
	internalLog "github.com/theTardigrade/runner/internal/log"
)

var (
	runMutex   sync.Mutex
	stopMutex  sync.Mutex
	exitMutex  sync.Mutex
	exited     bool
	ctx        context.Context
	cancelFunc context.CancelFunc
)

func daemon() {
	var args []string

	flag.Visit(func(f *flag.Flag) {
		if f.Name != "daemon" {
			args = append(args, "-"+f.Name+"="+f.Value.String())
		}
	})

	cmd := exec.Command(os.Args[0], args...)
	internalErrors.Check(cmd.Start())

	os.Exit(0)
}

func list() {
	files, err := ioutil.ReadDir(gobinPath)
	internalErrors.Check(err)

	var names []string

	for _, f := range files {
		name := f.Name()

		if isPathHidden(name) {
			continue
		}

		if strings.HasSuffix(name, pathWindowsNameSuffix) {
			name = strings.TrimSuffix(name, pathWindowsNameSuffix)
		}

		names = append(names, name)
	}

	{
		var b strings.Builder
		l := len(names)

		_, err = b.WriteString("FOUND %d COMMANDS")
		internalErrors.Check(err)

		if l > 0 {
			internalErrors.Check(b.WriteByte(':'))
		}

		internalFmt.Printf(b.String(), l)
	}

	for _, name := range names {
		internalFmt.Printf("%s[%s]", internalFmt.FourSpaces, name)
	}
}

func stop() {
	defer stopMutex.Unlock()
	stopMutex.Lock()

	if cancelFunc != nil {
		cancelFunc()
		cancelFunc = nil
		if ctx != nil {
			<-ctx.Done()
			ctx = nil
		}
	}
}

func run(path string) {
	defer runMutex.Unlock()
	runMutex.Lock()

	exitChan := make(chan struct{})

	func(c chan<- struct{}) {
		defer exitMutex.Unlock()
		exitMutex.Lock()

		if exited {
			c <- struct{}{}
		}
	}(exitChan)

	select {
	case <-exitChan:
		return
	default: // no-op
	}

	if *internalFlag.Log {
		internalLog.Open()
	}

	stopMutex.Lock()
	ctx, cancelFunc = context.WithCancel(context.Background())
	cmd := exec.CommandContext(ctx, path, internalFlag.SliceArguments...)
	stopMutex.Unlock()

	cmd.Stdout = os.Stdout

	if *internalFlag.Log {
		internalLog.SetWriter(&cmd.Stderr)
	} else {
		cmd.Stderr = os.Stderr
	}

	if *internalFlag.Verbose {
		internalFmt.Printf("RUNNING COMMAND [%s]", *internalFlag.Command)
	}

	err := cmd.Run()
	if *internalFlag.Log {
		if err != nil {
			internalLog.WriteString(err.Error())
		}

		internalLog.Close()
	}

	if *internalFlag.Verbose {
		internalFmt.Printf("COMPLETED COMMAND [%s] (%s)", *internalFlag.Command, internalErrors.Judge(err))
	}

	stopMutex.Lock()
	ctx, cancelFunc = nil, nil
	stopMutex.Unlock()
}

func command() {
	if *internalFlag.Iterations == 0 {
		panic(errZeroIterations)
	}

	path := filepath.Join(gobinPath, *internalFlag.Command)

	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		path += pathWindowsNameSuffix

		if info, err = os.Stat(path); os.IsNotExist(err) {
			panic(errCommandNotFound)
		}
	}

	if !info.Mode().IsRegular() {
		panic(errCommandNotRegularFile)
	}

	if isPathHidden(path) {
		panic(errCommandHiddenFile)
	}

	for i, j := *internalFlag.Iterations, 1; ; j++ {
		run(path)
		if i > 0 && j == i {
			break
		}
		time.Sleep(*internalFlag.Sleep)
	}
}

func exit() {
	exitMutex.Lock()
	runMutex.Lock()

	exited = true

	stop()

	if *internalFlag.Log {
		internalLog.Close()
	}

	if *internalFlag.CleanAll {
		internalLog.CleanAll()
	} else if *internalFlag.Clean {
		internalLog.Clean()
	}

	os.Exit(0)
}

func main() {
	if *internalFlag.Daemon {
		daemon()
	}

	if *internalFlag.List {
		list()
	}

	if *internalFlag.Command != "" {
		command()
	}

	exit()
}
