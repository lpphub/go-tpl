package logger

import (
	"context"
	"fmt"
	"time"
)

// Logger 日志接口
type Logger interface {
	// Log 核心日志方法
	Log(ctx context.Context, level Level, msg string, fields ...Field)

	// WithContext 将字段添加到 context（可选，某些实现可能不需要）
	WithContext(ctx context.Context, fields ...Field) context.Context

	WithCaller(skip int) Logger
}

var defaultLogger Logger

// Init 初始化日志
func Init(opts ...Option) error {
	cfg := &Config{
		Level: InfoLevel,
	}
	for _, opt := range opts {
		opt(cfg)
	}

	var err error
	defaultLogger, err = NewZeroLogger(cfg)
	return err
}

// Field 日志字段
type Field struct {
	Key   string
	Value interface{}
	Type  FieldType
}

type FieldType int

const (
	StringType FieldType = iota
	IntType
	Int64Type
	Uint64Type
	BoolType
	Float64Type
	DurationType
	TimeType
	ErrorType
	AnyType
)

// Level 日志级别
type Level int8

const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

func (l Level) String() string {
	switch l {
	case DebugLevel:
		return "debug"
	case InfoLevel:
		return "info"
	case WarnLevel:
		return "warn"
	case ErrorLevel:
		return "error"
	case FatalLevel:
		return "fatal"
	default:
		return "unknown"
	}
}

// SetLogger 设置自定义 logger
func SetLogger(l Logger) {
	defaultLogger = l
}

// GetLogger 获取当前 logger
func GetLogger() Logger {
	return defaultLogger
}

// Context 辅助方法

// WithField 在 context 中添加单个字段
func WithField(ctx context.Context, key string, value interface{}) context.Context {
	return defaultLogger.WithContext(ctx, Any(key, value))
}

// WithFields 在 context 中添加多个字段
func WithFields(ctx context.Context, fields ...Field) context.Context {
	return defaultLogger.WithContext(ctx, fields...)
}

// 包级别的日志方法

// Debug 记录 debug 级别日志
func Debug(ctx context.Context, msg string, fields ...Field) {
	defaultLogger.Log(ctx, DebugLevel, msg, fields...)
}

// Debugf 格式化 debug 日志
func Debugf(ctx context.Context, format string, args ...interface{}) {
	defaultLogger.Log(ctx, DebugLevel, fmt.Sprintf(format, args...))
}

// Info 记录 info 级别日志
func Info(ctx context.Context, msg string, fields ...Field) {
	defaultLogger.Log(ctx, InfoLevel, msg, fields...)
}

// Infof 格式化 info 日志
func Infof(ctx context.Context, format string, args ...interface{}) {
	defaultLogger.Log(ctx, InfoLevel, fmt.Sprintf(format, args...))

}

// Warn 记录 warn 级别日志
func Warn(ctx context.Context, msg string, fields ...Field) {
	defaultLogger.Log(ctx, WarnLevel, msg, fields...)
}

// Warnf 格式化 warn 日志
func Warnf(ctx context.Context, format string, args ...interface{}) {
	defaultLogger.Log(ctx, WarnLevel, fmt.Sprintf(format, args...))
}

// Error 记录 error 级别日志
func Error(ctx context.Context, msg string, fields ...Field) {
	defaultLogger.Log(ctx, ErrorLevel, msg, fields...)
}

// Errorf 格式化 error 日志
func Errorf(ctx context.Context, format string, args ...interface{}) {
	defaultLogger.Log(ctx, ErrorLevel, fmt.Sprintf(format, args...))
}

// Errw 记录带 error 的日志
func Errw(ctx context.Context, err error, fields ...Field) {
	defaultLogger.Log(ctx, ErrorLevel, err.Error(), fields...)
}

// Fatal 记录 fatal 级别日志
func Fatal(ctx context.Context, msg string, fields ...Field) {
	defaultLogger.Log(ctx, FatalLevel, msg, fields...)
}

// Fatalf 格式化 fatal 日志
func Fatalf(ctx context.Context, format string, args ...interface{}) {
	defaultLogger.Log(ctx, FatalLevel, fmt.Sprintf(format, args...))
}

// Field 构造函数

func Str(key, value string) Field {
	return Field{Key: key, Value: value, Type: StringType}
}

func Int(key string, value int) Field {
	return Field{Key: key, Value: value, Type: IntType}
}

func Int64(key string, value int64) Field {
	return Field{Key: key, Value: value, Type: Int64Type}
}

func Uint64(key string, value uint64) Field {
	return Field{Key: key, Value: value, Type: Uint64Type}
}

func Bool(key string, value bool) Field {
	return Field{Key: key, Value: value, Type: BoolType}
}

func Float64(key string, value float64) Field {
	return Field{Key: key, Value: value, Type: Float64Type}
}

func Duration(key string, value time.Duration) Field {
	return Field{Key: key, Value: value, Type: DurationType}
}

func Time(key string, value time.Time) Field {
	return Field{Key: key, Value: value, Type: TimeType}
}

func Err(err error) Field {
	return Field{Key: "error", Value: err, Type: ErrorType}
}

func Any(key string, value interface{}) Field {
	return Field{Key: key, Value: value, Type: AnyType}
}
