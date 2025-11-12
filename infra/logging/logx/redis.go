package logx

import (
	"context"
	"errors"
	"fmt"
	"go-tpl/infra/logging"
	"net"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisLogger 自定义Redis客户端日志记录器
type RedisLogger struct {
	logger logging.Logger
}

// NewRedisLogger 创建新的Redis日志记录器
func NewRedisLogger() *RedisLogger {
	return &RedisLogger{
		logger: logging.GetLogger().WithCaller(2),
	}
}

func (l *RedisLogger) DialHook(next redis.DialHook) redis.DialHook {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		start := time.Now()
		conn, err := next(ctx, network, addr)

		fields := []logging.Field{
			{Key: "event", Value: "redis_dial"},
			{Key: "address", Value: addr},
			{Key: "duration_ms", Value: time.Since(start).Milliseconds()},
		}

		if err != nil {
			l.logger.Write(logging.ErrorLevel, fmt.Sprintf("redis connected failed: %v", err), fields...)
		} else {
			l.logger.Write(logging.InfoLevel, "redis connected", fields...)
		}

		return conn, err
	}
}

// ProcessHook 实现命令处理钩子
func (l *RedisLogger) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {
		start := time.Now()
		err := next(ctx, cmd)
		elapsed := time.Since(start)

		// 构建命令字符串
		cmdStr := l.buildCommandString(cmd)

		// 添加字段
		fields := []logging.Field{
			{Key: "event", Value: "redis_command"},
			{Key: "command", Value: cmdStr},
			{Key: "duration_ms", Value: elapsed.Milliseconds()},
		}

		// 添加上下文字段
		fields = append(fields, logging.WithContext(ctx).Fields...)

		// 根据是否有错误确定日志级别和消息
		msg := "redis command success"
		level := logging.InfoLevel

		if err != nil && !errors.Is(err, redis.Nil) {
			msg = fmt.Sprintf("redis command failed: %v", err)
			level = logging.ErrorLevel
		}

		// 对于慢查询，使用警告级别
		if elapsed > 100*time.Millisecond && err == nil {
			msg = fmt.Sprintf("redis slow query detected (%v)", elapsed)
			level = logging.WarnLevel
		}

		l.logger.Write(level, msg, fields...)

		return err
	}
}

// ProcessPipelineHook 实现管道处理钩子
func (l *RedisLogger) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return func(ctx context.Context, cmds []redis.Cmder) error {
		start := time.Now()
		err := next(ctx, cmds)
		elapsed := time.Since(start)

		// 记录管道执行的整体信息
		fields := []logging.Field{
			{Key: "event", Value: "redis_pipeline"},
			{Key: "command_count", Value: len(cmds)},
			{Key: "duration_ms", Value: elapsed.Milliseconds()},
		}

		fields = append(fields, logging.WithContext(ctx).Fields...)

		if err != nil {
			l.logger.Write(logging.ErrorLevel, fmt.Sprintf("redis pipeline failed: %v", err), fields...)
		} else {
			l.logger.Write(logging.InfoLevel, "redis pipeline executed", fields...)
		}

		return err
	}
}

// buildCommandString 构建命令字符串（隐藏敏感信息）
func (l *RedisLogger) buildCommandString(cmd redis.Cmder) string {
	args := cmd.Args()
	if len(args) == 0 {
		return cmd.Name()
	}

	// 转换参数为字符串
	var argStrs []string
	for _, arg := range args {
		argStrs = append(argStrs, fmt.Sprintf("%v", arg))
	}

	// 隐藏敏感命令的参数
	cmdName := strings.ToUpper(cmd.Name())
	if l.isSensitiveCommand(cmdName) {
		if len(argStrs) > 1 {
			// 保留命令名，隐藏参数
			return fmt.Sprintf("%s ****", argStrs[0])
		}
	}

	return strings.Join(argStrs, " ")
}

// isSensitiveCommand 判断是否为敏感命令
func (l *RedisLogger) isSensitiveCommand(cmdName string) bool {
	sensitiveCommands := map[string]bool{
		"AUTH":    true,
		"SET":     false, // 可能包含敏感值
		"GET":     false,
		"HSET":    false, // 可能包含敏感值
		"HGET":    false,
		"LPUSH":   false, // 可能包含敏感值
		"RPUSH":   false, // 可能包含敏感值
		"LSET":    false, // 可能包含敏感值
		"CONFIG":  true,  // 可能包含敏感配置
		"DEBUG":   true,
		"EVAL":    true, // 脚本可能包含敏感信息
		"SCRIPT":  true,
		"MIGRATE": true,
		"RESTORE": true,
	}
	return sensitiveCommands[cmdName]
}
