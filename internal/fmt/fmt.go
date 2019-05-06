package fmt

import (
	"fmt"
	"runtime"
	"strings"

	internalTime "github.com/theTardigrade/runner/internal/time"
)

var (
	FourSpaces = strings.Repeat(" ", 4)
	Newline    string
)

func init() {
	if runtime.GOOS == "windows" {
		Newline = "\r\n"
	} else {
		Newline = "\n"
	}
}

func Print(s ...string) {
	fmt.Printf("%s[%d]%s%s\n", FourSpaces, internalTime.UnixMilli(), FourSpaces, strings.Join(s, " "))
}

func Printf(pattern string, s ...interface{}) {
	Print(fmt.Sprintf(pattern, s...))
}
