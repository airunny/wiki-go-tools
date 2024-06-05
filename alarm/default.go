package alarm

type LogLevel int

const (
	Debug LogLevel = iota
	Info
	Warn
	Error
)

type option struct {
	level LogLevel
}

type Option func(o *option)

func WithLogLevel(level LogLevel) Option {
	return func(o *option) {
		o.level = level
	}
}

var (
	defaultAlarm *Alarm
)

func NewDefaultAlarm(c *Config) (err error) {
	defaultAlarm, err = NewAlarm(c)
	return
}

func FeiShuAlarm(reqId string, msg string) {
	if defaultAlarm != nil {
		defaultAlarm.Alarm(reqId, msg)
	}
}

func FeiShuAlarmNow(reqId string, msg string, opts ...Option) {
	if defaultAlarm != nil {
		defaultAlarm.AlarmNow(reqId, msg, opts...)
	}
}

func Close() {
	if defaultAlarm != nil {
		defaultAlarm.Close()
	}
}
