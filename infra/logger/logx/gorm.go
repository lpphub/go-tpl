package logx

import (
	"context"
	"errors"
	"fmt"
	"go-tpl/infra/logger"
	"runtime"
	"strings"
	"time"

	"gorm.io/gorm"
	glog "gorm.io/gorm/logger"
)

// GormLogger 自定义GORM日志记录器
type GormLogger struct {
	logger        logger.Logger
	logLevel      glog.LogLevel
	slowThreshold time.Duration
}

// NewGormLogger 创建新的GORM日志记录器
func NewGormLogger() glog.Interface {
	return &GormLogger{
		logger:   logger.GetLogger().WithCaller(3),
		logLevel: glog.Info,
	}
}

// LogMode 设置日志模式
func (l *GormLogger) LogMode(level glog.LogLevel) glog.Interface {
	newLogger := *l
	newLogger.logLevel = level
	return &newLogger
}

// Info 记录信息日志
func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.logLevel >= glog.Info {
		l.log(ctx, glog.Info, msg, data...)
	}
}

// Warn 记录警告日志
func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.logLevel >= glog.Warn {
		l.log(ctx, glog.Warn, msg, data...)
	}
}

// Error 记录错误日志
func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.logLevel >= glog.Error {
		l.log(ctx, glog.Error, msg, data...)
	}
}

// Trace 记录SQL执行追踪
func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.logLevel <= glog.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	// add field
	fields := []logger.Field{
		{Key: "duration_ms", Value: elapsed.Milliseconds()},
		{Key: "sql", Value: sql},
		{Key: "rows", Value: rows},
	}

	msg := "sql do success"
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		msg = err.Error()
	}

	l.logger.Log(ctx, logger.InfoLevel, msg, fields...)
}

// log 通用日志记录方法
func (l *GormLogger) log(ctx context.Context, level glog.LogLevel, msg string, data ...interface{}) {
	if len(data) > 0 {
		msg = fmt.Sprintf("%s | data: %v", msg, data)
	}

	switch level {
	case glog.Info:
		l.logger.Log(ctx, logger.InfoLevel, msg)
	case glog.Warn:
		l.logger.Log(ctx, logger.WarnLevel, msg)
	case glog.Error:
		l.logger.Log(ctx, logger.ErrorLevel, msg)
	default:
		l.logger.Log(ctx, logger.InfoLevel, msg)
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
