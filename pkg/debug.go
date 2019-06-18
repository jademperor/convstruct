package pkg

import "fmt"

var (
	debug = true
)

// CloseDebug .
func CloseDebug() {
	debug = false
}

// DebugF .
func DebugF(format string, args ...interface{}) {
	if debug {
		fmt.Printf(format+"\n", args...)
	}
}
