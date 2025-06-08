package utils

import (
	"fmt"
	"path/filepath"
	"runtime"
	"time"
)

type logLevel string

const (
	LevelInfo  logLevel = "INFO"
	LevelError logLevel = "ERROR"
	LevelDebug logLevel = "DEBUG"
	LevelWarn  logLevel = "WARN"
)

type logger struct{}

var Logger = logger{}

func (l logger) Info(msg string, args ...interface{}) {
	l.log(LevelInfo, msg, args...)
}

func (l logger) Error(msg string, args ...interface{}) {
	l.log(LevelError, msg, args...)
}

func (l logger) Debug(msg string, args ...interface{}) {
	l.log(LevelDebug, msg, args...)
}

func (l logger) Warn(msg string, args ...interface{}) {
	l.log(LevelWarn, msg, args...)
}

func (l logger) log(level logLevel, msg string, args ...interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05.000")
	file, line := getCallerInfo(3) // 3 levels up to reach user call
	prefix := fmt.Sprintf("[%s] [%s] %s:%d", timestamp, level, file, line)
	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}
	fmt.Println(prefix, "-", msg)
}

func getCallerInfo(depth int) (string, int) {
	_, file, line, ok := runtime.Caller(depth)
	if !ok {
		return "???", 0
	}
	return filepath.Base(file), line
}
