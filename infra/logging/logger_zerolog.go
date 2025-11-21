package logging

import (
	"io"
	"os"
	"path/filepath"
	"strconv"

	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

type ZerologLogger struct {
	cfg    *Config
	logger zerolog.Logger
}

// Write 实现Logger接口
func (l *ZerologLogger) Write(level LogLevel, msg string, fields ...Field) {
	event := l.getEvent(level)

	// 添加字段
	for _, f := range fields {
		event = event.Interface(f.Key, f.Value)
	}

	event.CallerSkipFrame(3).Msg(msg)
}

func (l *ZerologLogger) WithCaller(skip int) Logger {
	return &ZerologLogger{
		cfg:    l.cfg,
		logger: l.logger.With().CallerWithSkipFrameCount(skip + 2).Logger(),
	}
}

func setupZeroLogger(cfg *Config) (Logger, error) {
	zl := &ZerologLogger{
		cfg: cfg,
	}
	zl.logger = zl.newLogger().With().Caller().Logger()
	return zl, nil
}

func (l *ZerologLogger) GetLogger() zerolog.Logger {
	return l.logger
}

func (l *ZerologLogger) newLogger() zerolog.Logger {
	var writer io.Writer
	if l.cfg.LogFile == "" {
		writer = zerolog.ConsoleWriter{Out: os.Stdout}
	} else {
		writer = l.getLogWriter(l.cfg.LogFile)
	}

	// 配置时间格式
	zerolog.TimeFieldFormat = "2006-01-02 15:04:05.999999"
	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		return filepath.Base(file) + ":" + strconv.Itoa(line)
	}

	return zerolog.New(writer).Level(l.getZerologLevel(l.cfg.LogLevel)).With().Timestamp().Logger()
}

func (l *ZerologLogger) getLogWriter(logFile string) io.Writer {
	lbj := &lumberjack.Logger{
		Filename:   logFile, // 日志文件名
		MaxSize:    300,     // 单个日志文件最大大小，单位 MB
		MaxBackups: 5,       // 保留的旧日志文件最大数量
		MaxAge:     14,      // 保留的旧日志文件最大天数
		Compress:   false,   // 是否压缩旧日志文件
	}

	return lbj
}

// getEvent 根据日志级别获取对应的事件
func (l *ZerologLogger) getEvent(level LogLevel) *zerolog.Event {
	switch level {
	case DebugLevel:
		return l.logger.Debug()
	case InfoLevel:
		return l.logger.Info()
	case WarnLevel:
		return l.logger.Warn()
	case ErrorLevel:
		return l.logger.Error()
	case FatalLevel:
		return l.logger.Fatal()
	default:
		return l.logger.Info()
	}
}

// getZerologLevel 转换日志级别
func (l *ZerologLogger) getZerologLevel(level LogLevel) zerolog.Level {
	switch level {
	case DebugLevel:
		return zerolog.DebugLevel
	case InfoLevel:
		return zerolog.InfoLevel
	case WarnLevel:
		return zerolog.WarnLevel
	case ErrorLevel:
		return zerolog.ErrorLevel
	case FatalLevel:
		return zerolog.FatalLevel
	default:
		return zerolog.InfoLevel
	}
}
