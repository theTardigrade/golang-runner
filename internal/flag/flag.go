package flag

import (
	"flag"
	"strings"
	"time"
)

var (
	Arguments  *string
	Clean      *bool
	CleanAll   *bool
	Command    *string
	Daemon     *bool
	Iterations *int
	List       *bool
	Log        *bool
	Sleep      *time.Duration
	Verbose    *bool

	SliceArguments []string
)

const (
	minSleepDuration = time.Nanosecond
)

func init() {
	Arguments = flag.String("args", "", "arguments to be supplied to command")
	Clean = flag.Bool("clean", false, "remove any temporary files that were created by the current invocation of the program")
	CleanAll = flag.Bool("clean-all", false, "remove any temporary files that are found, including those created by past invocations of the program")
	Command = flag.String("command", "", "name of command to execute")
	Daemon = flag.Bool("daemon", false, "run as a daemon")
	Iterations = flag.Int("iterations", -1, "maximum number of iterations; a negative value will loop infinitely")
	List = flag.Bool("list", false, "list all possible commands")
	Log = flag.Bool("log", false, "write errors to a temporary log file")
	Sleep = flag.Duration("sleep", minSleepDuration, "duration to sleep in between rerunning the command")
	Verbose = flag.Bool("verbose", false, "provide a greater level of output")

	flag.Parse()

	if *Sleep < minSleepDuration {
		*Sleep = minSleepDuration
	}

	if *Arguments != "" {
		for _, a := range strings.Split(*Arguments, " ") {
			if a != "" {
				SliceArguments = append(SliceArguments, a)
			}
		}
	}
}
