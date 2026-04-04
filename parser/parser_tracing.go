package parser

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

var traceLevel int = 0
var traceEnabled = readTraceEnabled()

const traceIdentPlaceholder string = "\t"

func readTraceEnabled() bool {
	enabled, err := strconv.ParseBool(os.Getenv("TRACE"))
	return err == nil && enabled
}

func identLevel() string {
	return strings.Repeat(traceIdentPlaceholder, traceLevel-1)
}

func tracePrint(fs string) {
	fmt.Printf("%s%s\n", identLevel(), fs)
}

func incIdent() { traceLevel = traceLevel + 1 }
func decIdent() { traceLevel = traceLevel - 1 }

func trace(msg string) string {
	if !traceEnabled {
		return msg
	}
	incIdent()
	tracePrint("BEGIN " + msg)
	return msg
}

func untrace(msg string) {
	if !traceEnabled {
		return
	}
	tracePrint("END " + msg)
	decIdent()
}
