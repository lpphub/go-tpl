package logger

type options struct {
	logFile string
}

type Option func(*options)

// Setup 初始化日志
func Setup(opts ...Option) {
	o := &options{}
	for _, opt := range opts {
		opt(o)
	}
	setupLogger(o.logFile)
}

// WithLogFile 设置日志文件
func WithLogFile(logFile string) Option {
	return func(o *options) {
		o.logFile = logFile
	}
}
