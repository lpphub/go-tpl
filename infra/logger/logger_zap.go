package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type options struct {
	logFile string
}

type Option func(*options)

func WithLogFile(logFile string) Option {
	return func(o *options) {
		o.logFile = logFile
	}
}

var logger *zap.Logger

func Setup(opts ...Option) {
	o := &options{}
	for _, opt := range opts {
		opt(o)
	}
	logger = newZapLogger(o.logFile)
}

func GetLogger() *zap.Logger {
	return logger
}

func Sync() {
	if logger != nil {
		_ = logger.Sync()
	}
}

func newZapLogger(logFile string) *zap.Logger {
	encoder := getLogEncoder()
	writer := getLogWriter(logFile)
	core := zapcore.NewCore(encoder, writer, zap.InfoLevel)
	return zap.New(core)
}

func getLogEncoder() zapcore.Encoder {
	encoderCfg := zapcore.EncoderConfig{
		LevelKey:       "level",
		TimeKey:        "time",
		CallerKey:      "file",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeCaller:   zapcore.ShortCallerEncoder, // 短路径编码器
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.999999"),
		EncodeDuration: zapcore.StringDurationEncoder,
	}
	return zapcore.NewJSONEncoder(encoderCfg)
}

func getLogWriter(logFile string) zapcore.WriteSyncer {
	if logFile == "" {
		return zapcore.AddSync(os.Stdout)
	}

	l := &lumberjack.Logger{
		Filename:   logFile, // 日志文件名
		MaxSize:    300,     // 单个日志文件最大大小，单位 MB
		MaxBackups: 5,       // 保留的旧日志文件最大数量
		MaxAge:     14,      // 保留的旧日志文件最大天数
		Compress:   false,   // 是否压缩旧日志文件
	}

	// 使用缓冲写入
	//return &zapcore.BufferedWriteSyncer{
	//	WS:            zapcore.AddSync(l),
	//	Size:          128 * 1024, // 128KB
	//	FlushInterval: 3 * time.Second,
	//}

	return zapcore.AddSync(l)
}
