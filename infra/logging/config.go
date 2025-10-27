package logging

type Config struct {
	provider LoggerProvider // 日志实现提供者
	LogLevel LogLevel       // 日志级别
	LogFile  string         // 日志文件 为空时输出到标准输出
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

// WithProvider 配置日志提供者
func WithProvider(provider LoggerProvider) Option {
	return func(o *Config) {
		o.provider = provider
	}
}
