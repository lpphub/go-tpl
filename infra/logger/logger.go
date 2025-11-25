package logger

import (
	"context"
	"io"
	"os"
)

// Level 日志级别
type Level int8

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
	FATAL
)

// Field 日志字段
type Field struct {
	Key   string
	Value any
}

// Logger 日志接口 - 只有两个核心方法
type Logger interface {
	// Log 输出日志
	Log(level Level, msg string, fields ...Field)
	Logc(callerSkip int, level Level, msg string, fields ...Field)

	// With 创建新的Logger
	With(fields ...Field) Logger
	// WithCaller 创建新的Logger, 设置调用栈跳过层数
	WithCaller(skip int) Logger
}

func New(opts ...Option) Logger {
	c := applyConfig(opts...)
	return newZapLogger(c)
}

// ==================== Context ====================

type ctxKey struct{}

func Ctx(ctx context.Context) Logger {
	if l, ok := ctx.Value(ctxKey{}).(Logger); ok {
		return l
	}
	return Default()
}

func ToCtx(ctx context.Context, l Logger) context.Context {
	return context.WithValue(ctx, ctxKey{}, l)
}

func WithCtx(ctx context.Context, fields ...Field) context.Context {
	return ToCtx(ctx, Ctx(ctx).With(fields...))
}

// ==================== 配置 ====================

type config struct {
	level  Level
	output io.Writer
}
type Option func(*config)

func applyConfig(opts ...Option) *config {
	c := &config{level: INFO, output: os.Stdout}
	for _, o := range opts {
		o(c)
	}
	return c
}

func WithLevel(l Level) Option      { return func(c *config) { c.level = l } }
func WithOutput(w io.Writer) Option { return func(c *config) { c.output = w } }

func Str(k, v string) Field             { return Field{k, v} }
func Int(k string, v int) Field         { return Field{k, v} }
func Int64(k string, v int64) Field     { return Field{k, v} }
func Float64(k string, v float64) Field { return Field{k, v} }
func Bool(k string, v bool) Field       { return Field{k, v} }
func Err(e error) Field                 { return Field{"error", e} }
func Any(k string, v any) Field         { return Field{k, v} }
