package logger

import "io"

type Config struct {
	Writer       io.Writer
	Level        Level
	EnableCaller bool
}

type Option func(*Config)

// WithWriter 设备日志输出
func WithWriter(writer io.Writer) Option {
	return func(o *Config) {
		o.Writer = writer
	}
}

// WithLogLevel 设置日志级别
func WithLogLevel(logLevel Level) Option {
	return func(o *Config) {
		o.Level = logLevel
	}
}
