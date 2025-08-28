package logext

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	logger2 "github.com/lpphub/golib/logger"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	ormutil "gorm.io/gorm/utils"
)

type GormLogger struct {
	Addr      string
	Database  string
	MaxSqlLen int
	logger    logger2.Logger
}

func NewGormLogger(db, addr string, sqlLen int) logger.Interface {
	if logger2.Log() == nil {
		return logger.Default
	}
	if sqlLen == 0 {
		sqlLen = 1024
	}
	return &GormLogger{
		Database:  db,
		Addr:      addr,
		MaxSqlLen: sqlLen,
		logger:    logger2.Log().With().CallerWithSkipFrameCount(5).Logger(),
	}
}

func (l *GormLogger) LogMode(logger.LogLevel) logger.Interface {
	return l
}

func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	m := fmt.Sprintf(msg, append([]interface{}{ormutil.FileWithLineNum()}, data...)...)
	// 非trace日志改为debug级别输出
	l.logger.Debug().Fields(l.commonFields(ctx)).Msg(m)
}

func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	m := fmt.Sprintf(msg, append([]interface{}{ormutil.FileWithLineNum()}, data...)...)
	l.logger.Warn().Fields(l.commonFields(ctx)).Msg(m)
}

func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	m := fmt.Sprintf(msg, append([]interface{}{ormutil.FileWithLineNum()}, data...)...)
	l.logger.Error().Fields(l.commonFields(ctx)).Msg(m)
}

func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	elapsed := time.Now().Sub(begin)
	cost := float64(elapsed.Nanoseconds()/1e4) / 100.0

	// 请求是否成功
	msg := "mysql do success"
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		// 没有找到记录不统计在请求错误中
		msg = err.Error()
	}

	sql, rows := fc()
	if l.MaxSqlLen <= 0 {
		sql = ""
	} else if len(sql) > l.MaxSqlLen {
		sql = sql[:l.MaxSqlLen]
	}

	l.logger.Info().Int64("affected_row", rows).
		Float64("cost_ms", cost).
		Str("sql", sql).
		Fields(l.commonFields(ctx)).
		Msg(msg)
}

func (l *GormLogger) commonFields(ctx context.Context) map[string]interface{} {
	var logId string
	if c, ok := ctx.(*gin.Context); ok && c != nil {
		logId, _ = ctx.Value("_ctx_log_id").(string)
	}
	fields := map[string]interface{}{
		"logId":   logId,
		"service": "mysql",
		"db":      l.Database,
		//"addr":    l.Addr,
	}
	return fields
}
