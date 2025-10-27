package logging

import "fmt"

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

var globalLogger Logger

func Init(opts ...Option) error {
	cfg := &Config{
		LogLevel: InfoLevel,
	}
	for _, opt := range opts {
		opt(cfg)
	}

	var err error
	globalLogger, err = setupZapLogger(cfg) // 使用zap作为默认日志库
	if err != nil {
		return fmt.Errorf("log: create logger failed: %v", err)
	}
	return nil
}
func GetLogger() Logger {
	return globalLogger
}
