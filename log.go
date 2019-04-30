package main

import (
	"fmt"
	"io/ioutil"
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
	logFileExt                          = ".log"
	logFilePattern                      = "runner-" + logFilePatternTimestampVariableName + "-x-*" + logFileExt
)

func openLogFile() {
	timestamp := strconv.FormatInt(unixMilli(), 10)

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

func cleanLogFiles() {
	pattern := strings.Replace(logFilePattern, logFilePatternTimestampVariableName, "[0-9]*", 1)
	pattern = strings.Replace(pattern, "*"+logFileExt, "[0-9]*"+logFileExt, 1)
	pattern = filepath.Join(os.TempDir(), pattern)

	matches, err := filepath.Glob(pattern)
	checkErr(err)

	for _, m := range matches {
		checkErr(os.Remove(m))

		if *flagVerbose {
			printf("DELETED FILE [%s]", m)
		}
	}
}
