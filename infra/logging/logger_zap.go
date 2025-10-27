package logging

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type ZapLogger struct {
	cfg    *Config
	logger *zap.Logger
}

// Write 实现Logger接口
func (l *ZapLogger) Write(level LogLevel, msg string, fields ...Field) {
	fieldSlice := make([]zap.Field, len(fields))
	for i, f := range fields {
		fieldSlice[i] = zap.Any(f.Key, f.Value)
	}

	switch level {
	case DebugLevel:
		l.logger.Debug(msg, fieldSlice...)
	case InfoLevel:
		l.logger.Info(msg, fieldSlice...)
	case WarnLevel:
		l.logger.Warn(msg, fieldSlice...)
	case ErrorLevel:
		l.logger.Error(msg, fieldSlice...)
	case FatalLevel:
		l.logger.Fatal(msg, fieldSlice...)
	default:
		l.logger.Info(msg, fieldSlice...)
	}
}

func setupZapLogger(cfg *Config) (Logger, error) {
	zl := &ZapLogger{
		cfg: cfg,
	}
	zl.logger = zl.newLogger()
	return zl, nil
}

func (l *ZapLogger) GetLogger() *zap.Logger {
	return l.logger
}

func (l *ZapLogger) newLogger() *zap.Logger {
	core := zapcore.NewCore(
		l.getLogEncoder(),
		l.getLogWriter(l.cfg.LogFile),
		l.getZapLevel(l.cfg.LogLevel),
	)
	return zap.New(core,
		zap.AddCaller(),
		zap.AddCallerSkip(3),
	)
}

func (l *ZapLogger) getLogEncoder() zapcore.Encoder {
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

func (l *ZapLogger) getLogWriter(logFile string) zapcore.WriteSyncer {
	if logFile == "" {
		return zapcore.AddSync(os.Stdout)
	}

	lbj := &lumberjack.Logger{
		Filename:   logFile, // 日志文件名
		MaxSize:    300,     // 单个日志文件最大大小，单位 MB
		MaxBackups: 5,       // 保留的旧日志文件最大数量
		MaxAge:     14,      // 保留的旧日志文件最大天数
		Compress:   false,   // 是否压缩旧日志文件
	}

	// 使用缓冲写入
	//return &zapcore.BufferedWriteSyncer{
	//	WS:            zapcore.AddSync(lbj),
	//	Size:          128 * 1024, // 128KB
	//	FlushInterval: 3 * time.Second,
	//}

	return zapcore.AddSync(lbj)
}

// getZapLevel 转换日志级别
func (l *ZapLogger) getZapLevel(level LogLevel) zapcore.Level {
	switch level {
	case DebugLevel:
		return zapcore.DebugLevel
	case InfoLevel:
		return zapcore.InfoLevel
	case WarnLevel:
		return zapcore.WarnLevel
	case ErrorLevel:
		return zapcore.ErrorLevel
	case FatalLevel:
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}
