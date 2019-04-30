package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var (
	logFile *os.File
)

const (
	logFilePatternTimestampVariableName = "$TIMESTAMP"
	logFilePattern                      = "runner-" + logFilePatternTimestampVariableName + "-x-*.log"
)

func openLogFile() {
	timestamp := strconv.FormatInt(int64(math.Round(float64(time.Now().UTC().UnixNano())/1e6)), 10) // milliseconds

	closeLogFile()

	var err error
	logFile, err = ioutil.TempFile("", strings.Replace(logFilePattern, logFilePatternTimestampVariableName, timestamp, 1))
	checkErr(err)

	_, err = logFile.WriteString(
		fmt.Sprintf(
			"DATE: %s%sCOMMAND: %s%sARGUMENTS: %s%s%s",
			time.Now().UTC().String(), newline, *flagCommand, newline, *flagArguments, newline, newline,
		),
	)
	checkErr(err)

	if *flagVerbose {
		printf("CREATED FILE [%s]", logFile.Name())
	}
}

func closeLogFile() {
	if logFile != nil {
		logFile.Close()
		logFile = nil
	}
}

var (
	cleanedLogFiles bool
)

func cleanLogFiles() {
	if cleanedLogFiles {
		return
	}

	pattern := filepath.Join(os.TempDir(), strings.Replace(logFilePattern, logFilePatternTimestampVariableName, "[0-9]*", 1))

	matches, err := filepath.Glob(pattern)
	checkErr(err)

	for _, m := range matches {
		checkErr(os.Remove(m))

		if *flagVerbose {
			printf("DELETED FILE [%s]", m)
		}
	}

	cleanedLogFiles = true
}
