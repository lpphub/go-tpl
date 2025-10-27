package logging

type Config struct {
	LogLevel LogLevel // 日志级别
	LogFile  string   // 日志文件 为空时输出到标准输出
}

type Option func(*Config)

// WithLogFile 设置日志文件
func WithLogFile(logFile string) Option {
	return func(o *Config) {
		o.LogFile = logFile
	}
}

// WithLogLevel 设置日志级别
func WithLogLevel(logLevel LogLevel) Option {
	return func(o *Config) {
		o.LogLevel = logLevel
	}
}
