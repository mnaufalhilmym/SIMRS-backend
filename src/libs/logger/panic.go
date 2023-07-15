package logger

import (
	"log"
	"os"
)

var plog = log.New(os.Stderr, "[PANIC]", 2)

func (l *Logger) Panic(message interface{}, option ...Option) {
	plog.Panicln("["+l.prefix+"]", message)
}
