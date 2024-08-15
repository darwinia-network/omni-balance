package logging

import (
	"fmt"
	"log"
	"os"
)

type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

var (
	currentLevel Level
	logger       *log.Logger
)

func init() {
	currentLevel = LevelDebug
	logger = log.New(os.Stdout, "", log.LstdFlags|log.Llongfile)
	logger.SetOutput(os.Stdout)
}

func (l Level) String() string {
	return [...]string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}[l]
}

func SetLevel(l Level) {
	currentLevel = l
	if currentLevel == LevelDebug {
		logger.SetFlags(log.LstdFlags | log.Llongfile)
	} else {
		logger.SetFlags(log.LstdFlags)
	}
}

func output(level Level, msg string, args ...interface{}) {
	if level >= currentLevel {
		logger.Output(3, fmt.Sprintf("[%s] %s", level.String(), fmt.Sprintf(msg, args...)))
	}
}

func Debugf(msg string, args ...interface{}) {
	output(LevelDebug, msg, args...)
}

func Debug(args ...interface{}) {
	output(LevelDebug, "", args...)
}

// Infof log
func Infof(msg string, args ...interface{}) {
	output(LevelInfo, msg, args...)
}

// Info log
func Info(args ...interface{}) {
	output(LevelInfo, "", args...)
}

// Errorf log
func Errorf(msg string, args ...interface{}) {
	output(LevelError, msg, args...)
}

// Error log
func Error(args ...interface{}) {
	output(LevelError, "", args...)
}

// Warnf log
func Warnf(msg string, args ...interface{}) {
	output(LevelWarn, msg, args...)
}

// Warn log
func Warn(args ...interface{}) {
	output(LevelWarn, "", args...)
}

// Fatalf send log fatalf
func Fatalf(msg string, args ...interface{}) {
	output(LevelFatal, msg, args...)
	os.Exit(1)
}

// Fatal send log fatal
func Fatal(args ...interface{}) {
	output(LevelFatal, "", args...)
	os.Exit(1)
}