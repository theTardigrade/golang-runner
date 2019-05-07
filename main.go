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
	internalStrings "github.com/theTardigrade/runner/internal/strings"
)

var (
	runMutex   sync.Mutex
	stopMutex  sync.Mutex
	exitMutex  sync.Mutex
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
		var b internalStrings.Builder
		l := len(names)

		b.WriteString("FOUND %d COMMANDS")

		if l > 0 {
			b.WriteByte(':')
		}

		internalFmt.Printf(b.String(), l)
	}

	for _, name := range names {
		internalFmt.Printf("%s[%s]", internalFmt.FourSpaces, name)
	}
}

func run(path string) {
	defer runMutex.Unlock()
	runMutex.Lock()

	waitIfExiting()

	if *internalFlag.Log {
		internalLog.Open()
	}

	stopMutex.Lock()
	ctx, cancelFunc = context.WithCancel(context.Background())
	cmd := exec.CommandContext(ctx, path, arguments...)
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
		waitIfExiting()
		time.Sleep(*internalFlag.Sleep)
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

func waitIfExiting() {
	exitMutex.Lock()
	exitMutex.Unlock()
}

func exit() {
	exitMutex.Lock()

	stop()

	runMutex.Lock()

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
