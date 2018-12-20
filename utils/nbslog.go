package utils

import (
	"github.com/op/go-logging"
	"os"
	"sync"
)

var instance *logging.Logger
var onceLog sync.Once

func GetLogInstance() *logging.Logger {

	onceLog.Do(func() {

		instance = newLogIns()
	})

	return instance
}

func newLogIns() *logging.Logger {

	log := logging.MustGetLogger("NBS")

	logFile, err := os.OpenFile(GetConfig().LogFileName,
		os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)

	if err != nil {
		panic(err)
	}

	fileBackend := logging.NewLogBackend(logFile, "-->", 0)

	fileFormat := logging.MustStringFormatter(
		`{time:01-02/15:04:05} %{longfunc:-30s} > %{level:.4s} %{message}`,
	)
	fileFormatBackend := logging.NewBackendFormatter(fileBackend, fileFormat)

	leveledFileBackend := logging.AddModuleLevel(fileFormatBackend)

	cmdFormat := logging.MustStringFormatter(
		`%{color}%{time:01-02/15:04:05} %{longfunc:-30s}> %{level:.4s} %{message}%{color:reset}`,
	)
	cmdBackend := logging.NewLogBackend(os.Stderr, ">>>", 0)
	formattedCmdBackend := logging.NewBackendFormatter(cmdBackend, cmdFormat)

	logging.SetBackend(leveledFileBackend, formattedCmdBackend)
	//logging.SetLevel(logging.DEBUG, "")

	return log
}

type Password string

func (p Password) Redacted() interface{} {
	return logging.Redact(string(p))
}

func SetLevel(level logging.Level, module string) {
	logging.SetLevel(level, module)
}

func Test() {

	log := GetLogInstance()

	log.Debugf("debug %s", Password("secret"))
	log.Info("info")
	log.Notice("notice")
	log.Warning("warning")
	log.Error("error")
	log.Critical("critical")
}
