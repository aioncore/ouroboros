package log

import (
	"fmt"
	"github.com/aioncore/ouroboros/pkg/service/utils"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"io"
	"os"
	"path/filepath"
	"time"
)

var (
	debugLogger log.Logger
	infoLogger  log.Logger
	errorLogger log.Logger
)

type Level string

const (
	DEBUG Level = "debug"
	INFO  Level = "info"
	ERROR Level = "error"
)

func InitLogger(filePath string, serviceType string) {
	debugLogger = log.NewLogfmtLogger(os.Stdout)
	debugLogger = level.NewFilter(debugLogger, level.AllowAll())

	infoLogFile, err := utils.OpenFile(filepath.Join(filePath, serviceType, "log", "info.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	infoLogger = log.NewLogfmtLogger(io.MultiWriter(os.Stdout, infoLogFile))
	infoLogger = level.NewFilter(infoLogger, level.AllowAll())

	errorLogFile, err := utils.OpenFile(filepath.Join(filePath, serviceType, "log", "error.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	errorLogger = log.NewLogfmtLogger(io.MultiWriter(os.Stdout, errorLogFile))
	errorLogger = level.NewFilter(errorLogger, level.AllowAll())
}

func Debug(format string, a ...any) {
	message := fmt.Sprintf(format, a...)
	logMessage(DEBUG, message)
}

func Info(format string, a ...any) {
	message := fmt.Sprintf(format, a...)
	logMessage(INFO, message)
}

func Error(format string, a ...any) {
	message := fmt.Sprintf(format, a...)
	logMessage(ERROR, message)
}

func logMessage(level Level, message string) {
	logger := levelLogger(level)
	err := logger.Log("time", time.Now(), "message", message)
	if err != nil {
		panic(err)
	}
}

func levelLogger(logLevel Level) log.Logger {
	switch logLevel {
	case DEBUG:
		return level.Debug(debugLogger)
	case INFO:
		return level.Info(infoLogger)
	case ERROR:
		return level.Error(errorLogger)
	default:
		return level.Info(infoLogger)
	}
}
