package logger

import (
	"fmt"
)

var logLevelsMap = map[string]int{
	"error": 1,
	"warn":  2,
	"info":  3,
	"debug": 4,
}

type Logger struct {
	level int
}

func New(level string) *Logger {
	logLevel, isListed := logLevelsMap[level]

	if !isListed {
		fmt.Printf("loglevel %s is not supported. supported levels\n: %v", level, logLevelsMap)
	}
	return &Logger{level: logLevel}
}

func (l Logger) Debug(msg string) {
	l.log(msg, logLevelsMap["debug"])
}

func (l Logger) Info(msg string) {
	l.log(msg, logLevelsMap["info"])
}

func (l Logger) Warn(msg string) {
	l.log(msg, logLevelsMap["warn"])

}

func (l Logger) Error(msg string) {
	l.log(msg, logLevelsMap["error"])

}

func (l Logger) log(msg string, logerLevel int) {
	if l.level >= logerLevel {
		fmt.Println(msg)
	}
}
