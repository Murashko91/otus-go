package logger

import (
	"fmt"
	"log"
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

func (l Logger) Debug(msg ...any) {
	l.log(logLevelsMap["debug"], msg...)
}

func (l Logger) Info(msg ...any) {
	l.log(logLevelsMap["info"], msg...)
}

func (l Logger) Warn(msg ...any) {
	l.log(logLevelsMap["warn"], msg...)
}

func (l Logger) Error(msg ...any) {
	l.log(logLevelsMap["error"], msg...)
}

func (l Logger) log(logerLevel int, msg ...any) {
	if l.level >= logerLevel {
		log.Println(msg...)
	}
}
