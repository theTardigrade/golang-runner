package main

import (
	"context"
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
		defer closeLogFile()
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
	if *flagLog && err != nil {
		logFile.WriteString(err.Error())
	}

	if *flagVerbose {
		printf("COMPLETED COMMAND [%s] (%s)", *flagCommand, judgeErr(err))
	}

	stopMutex.Lock()
	ctx, cancelFunc = nil, nil
	stopMutex.Unlock()
}

func list(path string) {
	files, err := ioutil.ReadDir(path)
	checkErr(err)

	var names []string

	for _, f := range files {
		name := f.Name()

		if strings.HasPrefix(name, pathHiddenNamePrefix) || strings.HasSuffix(name, pathHiddenNameSuffix) {
			continue
		}

		if isWindows {
			name = strings.TrimSuffix(name, pathWindowsNameSuffix)
		}

		names = append(names, name)
	}

	var b strings.Builder
	l := len(names)

	_, err = b.WriteString("FOUND %d COMMANDS")
	checkErr(err)

	if l > 0 {
		checkErr(b.WriteByte(':'))
	}

	printf(b.String(), l)

	for _, name := range names {
		printf("\t[%s]", name)
	}
}

func command(path string) {
	if *flagIterations == 0 {
		panic(errZeroIterations)
	}

	path = filepath.Join(path, *flagCommand)

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

	if !isWindows {
		result, err := exec.Command(filepath.Join(basePath, "isExecutable.sh"), path).Output()
		checkErr(err)
		if len(result) == 0 {
			panic(errCommandNotExecutable)
		}
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

	exited = true

	stop()

	runMutex.Lock()

	if *flagLog {
		closeLogFile()
	}

	if *flagClean {
		cleanLogFiles()
	}

	os.Exit(0)
}

func main() {
	path := gobin()

	if *flagList {
		list(path)
	}

	if *flagCommand != "" {
		command(path)
	}

	exit()
}
