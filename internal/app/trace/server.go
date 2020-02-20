package trace

import (
	"fmt"
	"os"

	"github.com/msmedes/scale/internal/pkg/trace"
)

var (
	tracePort = getEnv("TRACE", "5000")
	traceAddr = fmt.Sprintf("0.0.0.0:%s", tracePort)
)

// ServerListen creates a Trace Server
func ServerListen() {
	t := trace.NewTrace(traceAddr, tracePort)
	defer t.Shutdown()
	t.ServerListen()
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}
