package log

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	internalErrors "github.com/theTardigrade/runner/internal/errors"
	internalFlag "github.com/theTardigrade/runner/internal/flag"
	internalFmt "github.com/theTardigrade/runner/internal/fmt"
	internalTime "github.com/theTardigrade/runner/internal/time"
)

var (
	file  *os.File
	paths []string
)

const (
	filePatternTimestampVariableName = "$TIMESTAMP"
	fileExt                          = ".log"
	filePattern                      = "runner-" + filePatternTimestampVariableName + "-x-*" + fileExt
)

func Open() {
	timestamp := strconv.FormatInt(internalTime.UnixMilli(), 10)

	Close()

	var err error
	file, err = ioutil.TempFile("", strings.Replace(filePattern, filePatternTimestampVariableName, timestamp, 1))
	internalErrors.Check(err)

	filePath := file.Name()

	if *internalFlag.Clean || *internalFlag.CleanAll {
		paths = append(paths, filePath)
	}

	writeHeaders()

	if *internalFlag.Verbose {
		internalFmt.Printf("CREATED FILE [%s]", filePath)
	}
}

func writeHeaders() {
	dateHeader := "DATE: " + time.Now().UTC().String()
	commandHeader := "COMMAND: " + *internalFlag.Command
	argumentsHeader := "ARGUMENTS: " + *internalFlag.Arguments
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
	builder.WriteString(internalFmt.Newline)
	for _, h := range headers {
		builder.WriteString(h)
		builder.WriteString(internalFmt.Newline)
	}
	builder.WriteString(border)
	for i := 2; i > 0; i-- {
		builder.WriteString(internalFmt.Newline)
	}

	WriteString(builder.String())
}

func Close() {
	if file != nil {
		file.Close()
		file = nil
	}
}

func WriteString(s string) {
	_, err := file.WriteString(s)
	internalErrors.Check(err)
}

func SetWriter(w *io.Writer) {
	*w = file
}

func cleanOneCurried(identifier string) func(string) {
	return func(path string) {
		internalErrors.Check(os.Remove(path))

		if *internalFlag.Verbose {
			internalFmt.Printf("DELETED %s FILE [%s]", identifier, path)
		}
	}
}

func Clean() {
	cleanOne := cleanOneCurried("CURRENT")
	for _, p := range paths {
		cleanOne(p)
	}

	paths = []string{}
}

func CleanAll() {
	Clean()

	pattern := strings.Replace(filePattern, filePatternTimestampVariableName, "[0-9]*", 1)
	pattern = filepath.Join(os.TempDir(), pattern)

	matches, err := filepath.Glob(pattern)
	internalErrors.Check(err)

	cleanOne := cleanOneCurried("PREVIOUS")
	for _, m := range matches {
		cleanOne(m)
	}
}
