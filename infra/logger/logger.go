package logger

import (
	"time"

	"io"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func NewLogger(logFile string) *zap.Logger {
	// 控制台输出
	return zap.New(zapcore.NewCore(getLogEncoder(), getLogWriter(logFile), zap.InfoLevel))
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
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.999999"),
		EncodeDuration: zapcore.StringDurationEncoder,
	}
	return zapcore.NewJSONEncoder(encoderCfg)
}

func getLogWriter(logFile string) (ws zapcore.WriteSyncer) {
	var w io.Writer
	if logFile != "" {
		w = &lumberjack.Logger{
			Filename:   logFile,
			MaxSize:    400,
			MaxBackups: 5,
			MaxAge:     14,    // days
			Compress:   false, // disabled by default
		}
	} else {
		w = os.Stdout
	}

	//return zapcore.AddSync(w)

	// 开启缓冲区
	return &zapcore.BufferedWriteSyncer{
		WS:            zapcore.AddSync(w),
		Size:          128 * 1024, // 128KB
		FlushInterval: 3 * time.Second,
	}
}
