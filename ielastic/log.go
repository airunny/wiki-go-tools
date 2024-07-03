package ielastic

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/olivere/elastic/v7"
)

type LogLevel int

const (
	LogLevelInfo LogLevel = iota
	LogLevelError
	LogLevelTrace
)

type Log struct {
	level LogLevel
}

func (l *Log) Printf(format string, v ...interface{}) {
	switch l.level {
	case LogLevelError:
		log.Errorf(format, v...)
	case LogLevelTrace:
		log.Infof(format, v...)
	default:
		//log.Infof(format, v...)
	}
}

func NewLogWithLevel(level LogLevel) elastic.Logger {
	return &Log{
		level: level,
	}
}
