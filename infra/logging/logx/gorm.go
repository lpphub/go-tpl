package logx

import (
	"context"
	"errors"
	"fmt"
	"go-tpl/infra/logging"
	"runtime"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// GormLogger 自定义GORM日志记录器
type GormLogger struct {
	logger        logging.Logger
	logLevel      logger.LogLevel
	slowThreshold time.Duration
}

// NewGormLogger 创建新的GORM日志记录器
func NewGormLogger() logger.Interface {
	return &GormLogger{
		logger:   logging.GetLogger().WithCaller(1),
		logLevel: logger.Info,
	}
}

// LogMode 设置日志模式
func (l *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *l
	newLogger.logLevel = level
	return &newLogger
}

// Info 记录信息日志
func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.logLevel >= logger.Info {
		l.log(ctx, logger.Info, msg, data...)
	}
}

// Warn 记录警告日志
func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.logLevel >= logger.Warn {
		l.log(ctx, logger.Warn, msg, data...)
	}
}

// Error 记录错误日志
func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.logLevel >= logger.Error {
		l.log(ctx, logger.Error, msg, data...)
	}
}

// Trace 记录SQL执行追踪
func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.logLevel <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	// add field
	fields := []logging.Field{
		{Key: "duration_ms", Value: elapsed.Milliseconds()},
		{Key: "sql", Value: sql},
		{Key: "rows", Value: rows},
	}

	// 添加上下文字段
	fields = append(fields, logging.WithContext(ctx).GetFields()...)

	msg := "sql do success"
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		msg = err.Error()
	}

	l.logger.Write(logging.InfoLevel, msg, fields...)
}

// log 通用日志记录方法
func (l *GormLogger) log(_ context.Context, level logger.LogLevel, msg string, data ...interface{}) {
	if len(data) > 0 {
		msg = fmt.Sprintf("%s | data: %v", msg, data)
	}

	switch level {
	case logger.Info:
		l.logger.Write(logging.InfoLevel, msg)
	case logger.Warn:
		l.logger.Write(logging.WarnLevel, msg)
	case logger.Error:
		l.logger.Write(logging.ErrorLevel, msg)
	default:
		l.logger.Write(logging.InfoLevel, msg)
	}
}

// getCaller 获取调用者信息
func (l *GormLogger) getCaller() (string, int) {
	// 获取调用栈
	_, file, line, ok := runtime.Caller(4)
	if !ok {
		return "", 0
	}

	// 提取文件名（去掉路径）
	if idx := strings.LastIndex(file, "/"); idx >= 0 {
		file = file[idx+1:]
	}
	if idx := strings.LastIndex(file, "\\"); idx >= 0 {
		file = file[idx+1:]
	}

	return file, line
}
