package ilog

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"path"

	"github.com/airunny/wiki-go-tools/icontext"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	k8sJsonLogPath = "./jsonlog/%s"
)

type logger struct {
	*zap.SugaredLogger
}

type Helper struct {
	helper *log.Helper
	global *zap.SugaredLogger
}

func (h *Helper) Log(level log.Level, keyvals ...interface{}) error {
	h.helper.Log(level, keyvals...)
	return nil
}

func (h *Helper) Close() {
	_ = h.global.Sync()
}

func NewLogger(id, name string, opts ...Option) *Helper {
	o := &options{}
	for _, opt := range opts {
		opt(o)
	}

	var (
		logPath   = os.Getenv("LOG_PATH")
		_, k8sEnv = os.LookupEnv("KUBERNETES_SERVICE_HOST")
		writer    io.Writer
	)

	if logPath == "" && k8sEnv {
		logPath = fmt.Sprintf(k8sJsonLogPath, name)
	}

	if logPath == "" || o.console {
		writer = os.Stdout
	} else {
		writer = &lumberjack.Logger{
			Filename:  path.Join(logPath, fmt.Sprintf("%s-%s.log", name, getLocalIP())),
			MaxSize:   100,
			MaxAge:    7,
			LocalTime: true,
		}
	}

	conf := zap.NewProductionEncoderConfig()
	conf.EncodeTime = zapcore.ISO8601TimeEncoder
	conf.EncodeLevel = zapcore.CapitalLevelEncoder
	encoder := zapcore.NewJSONEncoder(conf)

	var (
		core         = zapcore.NewCore(encoder, zapcore.AddSync(writer), zapcore.DebugLevel)
		globalZapLog = zap.New(
			core,
			//zap.AddCaller(),
			zap.AddStacktrace(zap.ErrorLevel),
			//zap.AddCallerSkip(0),
		)
	)

	var (
		globalSugarLogger = globalZapLog.Sugar()
		kvs               = []interface{}{
			"service_id", id,
			"service_name", name,
			"trace_id", tracing.TraceID(),
			"span_id", tracing.SpanID(),
			//"caller", log.Caller(6),
		}
	)
	kvs = append(kvs, icontext.LoggerValues()...)

	ll := log.With(&logger{
		SugaredLogger: globalSugarLogger,
	}, kvs...)
	log.SetLogger(ll)

	return &Helper{
		helper: log.Context(context.Background()),
		global: globalSugarLogger,
	}
}

func (l *logger) Log(level log.Level, keyvals ...interface{}) error {
	if len(keyvals) == 0 || len(keyvals)%2 != 0 {
		return nil
	}

	var (
		msg string
		kvs []interface{}
	)

	for i := 0; i < len(keyvals); i += 2 {
		var (
			key   = fmt.Sprint(keyvals[i])
			value = fmt.Sprint(keyvals[i+1])
		)

		if value == "" {
			fmt.Println("这里的内容为空：", key, value)
			continue
		}

		if key == log.DefaultMessageKey {
			msg = value
			continue
		}

		kvs = append(kvs, zap.Any(key, value))
	}

	switch level {
	case log.LevelDebug:
		l.Debugw(msg, kvs...)
	case log.LevelInfo:
		l.Infow(msg, kvs...)
	case log.LevelWarn:
		l.Warnw(msg, kvs...)
	case log.LevelError:
		l.Errorw(msg, kvs...)
	case log.LevelFatal:
		l.Fatalw(msg, kvs...)
	}
	return nil
}

func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return uuid.New().String()
	}

	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String()
			}
		}
	}

	return uuid.New().String()
}
