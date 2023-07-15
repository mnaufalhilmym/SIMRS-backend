package logger

type Logger struct {
	prefix string
}

func New(prefix string) Logger {
	return Logger{
		prefix: prefix,
	}
}
