package scale

import (
	"strings"

	"go.uber.org/zap"
)

var traceLevel int

const traceIdentPlaceholder string = "\t"

func indentLevel() string {
	return strings.Repeat(traceIdentPlaceholder, traceLevel-1)
}

func tracePrint(fs string, logger *zap.SugaredLogger) {
	logger.Infof("%s%s\n", indentLevel(), fs)
}

func incrementIndex() { traceLevel = traceLevel + 1 }
func decrementIndex() { traceLevel = traceLevel - 1 }

// Trace starts the, well, trace
func Trace(msg string, logger *zap.SugaredLogger) string {
	incrementIndex()
	tracePrint("BEGIN "+msg, logger)
	return msg
}

// Untrace stops the trace
func Untrace(msg string, logger *zap.SugaredLogger) {
	tracePrint("END"+msg, logger)
	decrementIndex()
}
