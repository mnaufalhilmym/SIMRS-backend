package logger

import (
	"log"
	"os"
	"runtime/debug"
)

var nlog = log.New(os.Stdout, "[LOG]", 0)

func (l *Logger) Log(message interface{}, option ...*Option) {
	nlog.Println("["+l.prefix+"]", message)

	if len(option) == 0 {
		return
	}

	if option[0].IsPrintStack {
		debug.PrintStack()
	}
}
