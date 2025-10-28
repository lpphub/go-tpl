package logx

import (
	"context"
	"errors"
	"fmt"
	"go-tpl/infra/logging"
	"runtime"
	"strings"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// GormLogger 自定义GORM日志记录器
type GormLogger struct {
	logger        *zap.Logger
	LogLevel      logger.LogLevel
	slowThreshold time.Duration
}

// NewGormLogger 创建新的GORM日志记录器
func NewGormLogger() logger.Interface {
	zapLog := logging.GetLogger().(*logging.ZapLogger).GetLogger()
	return &GormLogger{
		logger:   zapLog.WithOptions(zap.AddCallerSkip(0)),
		LogLevel: logger.Info,
	}
}

// LogMode 设置日志模式
func (l *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *l
	newLogger.LogLevel = level
	return &newLogger
}

// Info 记录信息日志
func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Info {
		l.log(ctx, logger.Info, msg, data...)
	}
}

// Warn 记录警告日志
func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Warn {
		l.log(ctx, logger.Warn, msg, data...)
	}
}

// Error 记录错误日志
func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Error {
		l.log(ctx, logger.Error, msg, data...)
	}
}

// Trace 记录SQL执行追踪
func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Now().Sub(begin)
	duration := float64(elapsed.Nanoseconds()/1e4) / 100.0

	msg := "sql do success"
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		// 没有找到记录不统计在请求错误中
		msg = err.Error()
	}

	// 获取sql
	sql, rows := fc()

	fields := []zap.Field{
		zap.String("sql", sql),
		zap.Int64("rows", rows),
		zap.Float64("duration", duration),
	}
	// 获取上下文中的字段
	logCtx := logging.WithContext(ctx)
	for _, f := range logCtx.Fields {
		fields = append(fields, zap.Any(f.Key, f.Value))
	}

	l.logger.Info(msg, fields...)
}

// log 通用日志记录方法
func (l *GormLogger) log(ctx context.Context, level logger.LogLevel, msg string, data ...interface{}) {
	if len(data) > 0 {
		msg = fmt.Sprintf("%s | data: %v", msg, data)
	}

	switch level {
	case logger.Info:
		logging.Info(ctx, msg)
	case logger.Warn:
		logging.Warn(ctx, msg)
	case logger.Error:
		logging.Error(ctx, msg)
	default:
		logging.Info(ctx, msg)
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
