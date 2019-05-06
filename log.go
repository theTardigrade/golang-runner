package main

import (
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

	writeHeadersToLogFile()

	if *flagVerbose {
		printf("CREATED FILE [%s]", logFilePath)
	}
}

func writeHeadersToLogFile() {
	dateHeader := "DATE: " + time.Now().UTC().String()
	commandHeader := "COMMAND: " + *flagCommand
	argumentsHeader := "ARGUMENTS: " + *flagArguments
	headers := []string{dateHeader, commandHeader, argumentsHeader}

	var l int
	for _, h := range headers {
		if m := len(h); m > l {
			l = m
		}
	}

	border := strings.Repeat("*", l)

	var builder strings.Builder

	builder.WriteString(border)
	builder.WriteString(newline)
	for _, h := range headers {
		builder.WriteString(h)
		builder.WriteString(newline)
	}
	builder.WriteString(border)
	for i := 2; i > 0; i-- {
		builder.WriteString(newline)
	}

	_, err := logFile.WriteString(builder.String())
	checkErr(err)
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
