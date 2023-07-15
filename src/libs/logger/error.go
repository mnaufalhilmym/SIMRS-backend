package logger

import (
	"log"
	"os"
	"runtime/debug"
)

var elog = log.New(os.Stderr, "[ERROR]", 1)

func (l *Logger) Error(message interface{}, option ...*Option) {
	elog.Println("["+l.prefix+"]", message)

	if option == nil || option[0].IsPrintStack {
		debug.PrintStack()
	}

	if len(option) == 0 {
		return
	}

	if option[0].IsExit {
		exitCode := 1
		if option[0].ExitCode > 1 {
			exitCode = option[0].ExitCode
		}

		os.Exit(exitCode)
	}
}
