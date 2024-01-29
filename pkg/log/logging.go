package log

import (
	"path/filepath"
	"strings"
)

type Configuration struct {
	LogType  string
	LogFile  string
	LogLevel string

	RotateMaxSize    int
	RotateMaxAge     int
	RotateMaxBackups int
	Compress         bool
}

type LoggerInterface interface {
	Info(args ...interface{})
	Infof(f string, args ...interface{})
	Error(args ...interface{})
	Errorf(f string, args ...interface{})
	Warn(args ...interface{})
	Warnf(f string, args ...interface{})
}

var (
	Logger    LoggerInterface
	AccessLog LoggerInterface
)

func Register(logType, logDir, logLevel string) {
	// 支持 INFO, WARN 和 ERROR，默认为 INFO
	Level := "info"
	if strings.ToLower(logLevel) == "error" {
		Level = "error"
	} else if strings.ToLower(logLevel) == "warn" {
		Level = "warn"
	}

	AccessLog, _ = NewZapLogger(Configuration{
		LogType:          logType,
		LogFile:          filepath.Join(logDir, "access.log"), // 使用文件类型时生效
		LogLevel:         "info",                              // access 的 log 只会有 info
		RotateMaxSize:    500,
		RotateMaxAge:     7,
		RotateMaxBackups: 3,
	})

	Logger, _ = NewZapLogger(Configuration{
		LogType:          logType,
		LogFile:          filepath.Join(logDir, "oats.log"), // 使用文件类型时生效
		LogLevel:         Level,
		RotateMaxSize:    500,
		RotateMaxAge:     7,
		RotateMaxBackups: 3,
	})
}
