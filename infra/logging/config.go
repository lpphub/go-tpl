package logging

type config struct {
	logFile  string
	logLevel LogLevel
}

type Option func(*config)

// WithLogFile 设置日志文件
func WithLogFile(logFile string) Option {
	return func(o *config) {
		o.logFile = logFile
	}
}

// WithLogLevel 设置日志级别
func WithLogLevel(logLevel LogLevel) Option {
	return func(o *config) {
		o.logLevel = logLevel
	}
}
