package zlog

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Field 日志字段
type Field struct {
	Key   string
	Value interface{}
}

// Level 日志级别
type Level int8

const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
	PanicLevel
)

// Config 日志配置
type Config struct {
	OutputPath string
	Level      Level
}

var globalLogger zerolog.Logger

// Init 初始化日志
func Init(cfg *Config) {
	if cfg == nil {
		cfg = &Config{
			Level: InfoLevel,
		}
	}

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMicro
	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		return filepath.Base(file) + ":" + strconv.Itoa(line)
	}

	writer := getLogWriter(cfg.OutputPath)

	globalLogger = zerolog.New(writer).
		With().
		Timestamp().
		CallerWithSkipFrameCount(3).
		Logger().
		Level(getLevel(cfg.Level))
}

// getLevel 转换日志级别
func getLevel(level Level) zerolog.Level {
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
	case PanicLevel:
		return zerolog.PanicLevel
	default:
		return zerolog.InfoLevel
	}
}

func getLogWriter(logFile string) io.Writer {
	if logFile == "" {
		return os.Stdout
	}

	lbj := &lumberjack.Logger{
		Filename:   logFile, // 日志文件名
		MaxSize:    300,     // 单个日志文件最大大小，单位 MB
		MaxBackups: 5,       // 保留的旧日志文件最大数量
		MaxAge:     14,      // 保留的旧日志文件最大天数
		Compress:   false,   // 是否压缩旧日志文件
	}
	return lbj
}

func CallerWithSkip(skip int) zerolog.Logger {
	return globalLogger.With().CallerWithSkipFrameCount(skip).Logger()
}

func WithContext(ctx context.Context) context.Context {
	return globalLogger.WithContext(ctx)
}

// WithField 在 context 的 logger 中添加自定义字段
func WithField(ctx context.Context, key string, value interface{}) context.Context {
	return WithFields(ctx, Any(key, value))
}

// WithFields 在 context 的 logger 中添加多个字段
func WithFields(ctx context.Context, fields ...Field) context.Context {
	logger := zerolog.Ctx(ctx)
	if logger.GetLevel() == zerolog.Disabled {
		logger = &globalLogger
	}

	logCtx := logger.With()
	for _, field := range fields {
		logCtx = logCtx.Interface(field.Key, field.Value)
	}

	l := logCtx.Logger()
	return l.WithContext(ctx)
}

// getLogger 从 context 获取 logger，如果没有则返回 global logger
func getLogger(ctx context.Context) *zerolog.Logger {
	if ctx == nil {
		return &globalLogger
	}

	logger := zerolog.Ctx(ctx)
	// 如果 context 中没有 logger，zerolog.Ctx 返回 disabled logger
	if logger.GetLevel() == zerolog.Disabled {
		return &globalLogger
	}

	return logger
}

// applyFields 应用字段到事件
func applyFields(event *zerolog.Event, fields []Field) *zerolog.Event {
	for _, field := range fields {
		event = event.Interface(field.Key, field.Value)
	}
	return event
}

// Debug 级别日志

func Debug(ctx context.Context, msg string, fields ...Field) {
	event := getLogger(ctx).Debug()
	applyFields(event, fields).Msg(msg)
}

func Debugf(ctx context.Context, format string, args ...interface{}) {
	getLogger(ctx).Debug().Msgf(format, args...)
}

// Info 级别日志

func Info(ctx context.Context, msg string, fields ...Field) {
	event := getLogger(ctx).Info()
	applyFields(event, fields).Msg(msg)
}

func Infof(ctx context.Context, format string, args ...interface{}) {
	getLogger(ctx).Info().Msgf(format, args...)
}

// Warn 级别日志

func Warn(ctx context.Context, msg string, fields ...Field) {
	event := getLogger(ctx).Warn()
	applyFields(event, fields).Msg(msg)
}

func Warnf(ctx context.Context, format string, args ...interface{}) {
	getLogger(ctx).Warn().Msgf(format, args...)
}

// Error 级别日志

func Error(ctx context.Context, msg string, fields ...Field) {
	event := getLogger(ctx).Error()
	applyFields(event, fields).Msg(msg)
}

func Errorf(ctx context.Context, format string, args ...interface{}) {
	getLogger(ctx).Error().Msgf(format, args...)
}

func ErrorWithErr(ctx context.Context, err error, msg string, fields ...Field) {
	event := getLogger(ctx).Error().Err(err)
	applyFields(event, fields).Msg(msg)
}

// Fatal 级别日志

func Fatal(ctx context.Context, msg string, fields ...Field) {
	event := getLogger(ctx).Fatal()
	applyFields(event, fields).Msg(msg)
}

func Fatalf(ctx context.Context, format string, args ...interface{}) {
	getLogger(ctx).Fatal().Msgf(format, args...)
}

// Panic 级别日志

func Panic(ctx context.Context, msg string, fields ...Field) {
	event := getLogger(ctx).Panic()
	applyFields(event, fields).Msg(msg)
}

func Panicf(ctx context.Context, format string, args ...interface{}) {
	getLogger(ctx).Panic().Msgf(format, args...)
}

// Field 构造函数

func Str(key, value string) Field {
	return Field{Key: key, Value: value}
}

func Int(key string, value int) Field {
	return Field{Key: key, Value: value}
}

func Int64(key string, value int64) Field {
	return Field{Key: key, Value: value}
}

func Uint(key string, value uint) Field {
	return Field{Key: key, Value: value}
}

func Bool(key string, value bool) Field {
	return Field{Key: key, Value: value}
}

func Float64(key string, value float64) Field {
	return Field{Key: key, Value: value}
}

func Err(err error) Field {
	return Field{Key: "error", Value: err.Error()}
}

func Any(key string, value interface{}) Field {
	return Field{Key: key, Value: value}
}

func Duration(key string, value time.Duration) Field {
	return Field{Key: key, Value: value}
}

func Time(key string, value time.Time) Field {
	return Field{Key: key, Value: value}
}
