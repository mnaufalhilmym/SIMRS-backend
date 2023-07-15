package exithandler

import (
	"os"
	"os/signal"
	applogger "sim-puskesmas/src/libs/logger"
	"syscall"
)

var fnsRunInExit []FnRunInExit
var logger = applogger.New("ExitHandler")

type FnRunInExit struct {
	FnDescription string
	Fn            func()
}

func Add(newFns ...FnRunInExit) {
	fnsRunInExit = append(fnsRunInExit, newFns...)
}

type Config struct {
	IsDebug bool
}

func New(config ...*Config) {
	if len(config) == 0 {
		config = make([]*Config, 1)
	}

	if config[0].IsDebug {
		logger.Log("listen to exit signals")
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGTERM)
	go func() {
		<-c
		if config[0].IsDebug && len(fnsRunInExit) > 0 {
			logger.Log("clearing resources")
		}
		for _, fn := range fnsRunInExit {
			if config[0].IsDebug {
				logger.Log(fn.FnDescription)
			}
			fn.Fn()
		}
	}()
}
