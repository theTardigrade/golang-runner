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
	checkErr(cmd.Start())

	os.Exit(0)
}

func list() {
	files, err := ioutil.ReadDir(gobinPath)
	checkErr(err)

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
		checkErr(err)

		if l > 0 {
			checkErr(b.WriteByte(':'))
		}

		printf(b.String(), l)
	}

	for _, name := range names {
		printf("%s[%s]", fourSpaces, name)
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

	defer runMutex.Unlock()
	runMutex.Lock()

	if *flagLog {
		openLogFile()
	}

	stopMutex.Lock()
	ctx, cancelFunc = context.WithCancel(context.Background())
	cmd := exec.CommandContext(ctx, path, arguments...)
	stopMutex.Unlock()

	cmd.Stdout = os.Stdout

	if *flagLog {
		cmd.Stderr = logFile
	} else {
		cmd.Stderr = os.Stderr
	}

	if *flagVerbose {
		printf("RUNNING COMMAND [%s]", *flagCommand)
	}

	err := cmd.Run()
	if *flagLog {
		if err != nil {
			logFile.WriteString(err.Error())
		}
		closeLogFile()
	}

	if *flagVerbose {
		printf("COMPLETED COMMAND [%s] (%s)", *flagCommand, judgeErr(err))
	}

	stopMutex.Lock()
	ctx, cancelFunc = nil, nil
	stopMutex.Unlock()
}

func command() {
	if *flagIterations == 0 {
		panic(errZeroIterations)
	}

	path := filepath.Join(gobinPath, *flagCommand)

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

	for i, j := *flagIterations, 1; ; j++ {
		run(path)
		if i > 0 && j == i {
			break
		}
		time.Sleep(*flagSleep)
	}
}

func exit() {
	exitMutex.Lock()
	runMutex.Lock()

	exited = true

	stop()

	if *flagLog {
		closeLogFile()
	}

	if *flagCleanAll {
		cleanAllLogFiles()
	} else if *flagClean {
		cleanLogFiles()
	}

	os.Exit(0)
}

func main() {
	if *flagDaemon {
		daemon()
	}

	if *flagList {
		list()
	}

	if *flagCommand != "" {
		command()
	}

	exit()
}
