package ikafka

import "github.com/go-kratos/kratos/v2/log"

type Log struct{}

func (l *Log) Print(v ...interface{}) {
	log.Info(v...)
}

func (l *Log) Printf(format string, v ...interface{}) {
	log.Infof(format, v...)
}

func (l *Log) Println(v ...interface{}) {
	log.Info(v...)
}
