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
	logFile      *os.File
	logFilePaths []string
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

	logFilePath := logFile.Name()

	if *flagClean {
		logFilePaths = append(logFilePaths, logFilePath)
	}

	_, err = logFile.WriteString(
		fmt.Sprintf(
			"DATE: %s%sCOMMAND: %s%sARGUMENTS: %s%s%s",
			time.Now().UTC().String(), newline, *flagCommand, newline, *flagArguments, newline, newline,
		),
	)
	checkErr(err)

	if *flagVerbose {
		printf("CREATED FILE [%s]", logFilePath)
	}
}

func closeLogFile() {
	if logFile != nil {
		logFile.Close()
		logFile = nil
	}
}

func cleanLogFile(path string) {
	checkErr(os.Remove(path))

	if *flagVerbose {
		printf("DELETED FILE [%s]", path)
	}
}

func cleanLogFiles() {
	for _, p := range logFilePaths {
		cleanLogFile(p)
	}

	logFilePaths = []string{}
}

func cleanAllLogFiles() {
	pattern := strings.Replace(logFilePattern, logFilePatternTimestampVariableName, "[0-9]*", 1)
	pattern = filepath.Join(os.TempDir(), pattern)

	matches, err := filepath.Glob(pattern)
	checkErr(err)

	for _, m := range matches {
		cleanLogFile(m)
	}
}
