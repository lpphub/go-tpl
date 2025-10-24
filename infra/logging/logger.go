package logging

type LogLevel int8

const (
	DebugLevel LogLevel = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

type Field struct {
	Key   string
	Value interface{}
}

type Logger interface {
	Write(level LogLevel, msg string, fields ...Field)
}

var logger Logger

func Init(opts ...Option) {
	cfg := &config{}
	for _, opt := range opts {
		opt(cfg)
	}

	logger = setupZapLogger(cfg)
}
func GetLogger() Logger {
	return logger
}
