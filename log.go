package yiyun

import (
	"os"

	logging "github.com/op/go-logging"
)

var log = logging.MustGetLogger("example")
var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
)

func init() {
	backend1 := logging.NewLogBackend(os.Stderr, "", 0)
	backend2 := logging.NewLogBackend(os.Stderr, "", 0)

	// For messages written to backend2 we want to add some additional
	// information to the output, including the used log level and the name of
	// the function.
	backend2Formatter := logging.NewBackendFormatter(backend2, format)

	// Only errors and more severe messages should be sent to backend1
	backend1Leveled := logging.AddModuleLevel(backend1)
	backend1Leveled.SetLevel(logging.ERROR, "")

	// Set the backends to be used.
	logging.SetBackend(backend1Leveled, backend2Formatter)
}

//Info  Info日志
func Info(arg ...interface{}) {
	log.Info(arg)
}

//Error Error日志
func Error(arg ...interface{}) {
	log.Error(arg)
}

//Debug Debug日志
func Debug(arg ...interface{}) {
	log.Debug(arg)
}

//Panic Panic日志
func Panic(arg ...interface{}) {
	log.Panic(arg)
}
