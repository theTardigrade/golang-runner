package main

import (
	"flag"
	"strings"
	"time"
)

var (
	flagArguments  *string
	flagClean      *bool
	flagCommand    *string
	flagDaemon     *bool
	flagIterations *int
	flagList       *bool
	flagLog        *bool
	flagSleep      *time.Duration
	flagVerbose    *bool

	arguments []string
)

func init() {
	flagArguments = flag.String("args", "", "arguments to be supplied to command")
	flagClean = flag.Bool("clean", false, "remove old temporary files")
	flagCommand = flag.String("command", "", "name of command to execute")
	flagDaemon = flag.Bool("daemon", false, "run as a daemon")
	flagIterations = flag.Int("iterations", -1, "maximum number of iterations; a negative value will loop infinitely")
	flagList = flag.Bool("list", false, "list all possible commands")
	flagLog = flag.Bool("log", false, "write errors to a temporary log file")
	flagSleep = flag.Duration("sleep", time.Nanosecond, "duration to sleep in between rerunning the command")
	flagVerbose = flag.Bool("verbose", false, "provide a greater level of output")

	flag.Parse()

	if *flagSleep < time.Nanosecond {
		*flagSleep = time.Nanosecond
	}

	if *flagArguments != "" {
		for _, a := range strings.Split(*flagArguments, " ") {
			if a != "" {
				arguments = append(arguments, a)
			}
		}
	}
}
