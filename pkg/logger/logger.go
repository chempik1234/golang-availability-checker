package logger

import (
	"context"
	"go.uber.org/zap"
)

const (
	LoggerKey = "logger"
	RequestID = "request_id"
)

var logger *zap.Logger

type Logger struct {
	l *zap.Logger
}

func TryAppendRequestIDFromContext(ctx context.Context, fields []zap.Field) []zap.Field {
	if ctx.Value(RequestID) != nil {
		fields = append(fields, zap.String(RequestID, ctx.Value(RequestID).(string)))
	}
	return fields
}

func New(ctx context.Context) (context.Context, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}

	ctx = context.WithValue(ctx, LoggerKey, &Logger{l: logger})
	return ctx, nil
}

func GetLoggerFromCtx(ctx context.Context) *Logger {
	return ctx.Value(LoggerKey).(*Logger)
}

func (l *Logger) Info(ctx context.Context, msg string, fields ...zap.Field) {
	fields = TryAppendRequestIDFromContext(ctx, fields)
	l.l.Info(msg, fields...)
}

func (l *Logger) Warn(ctx context.Context, msg string, fields ...zap.Field) {
	fields = TryAppendRequestIDFromContext(ctx, fields)
	l.l.Warn(msg, fields...)
}

func (l *Logger) Error(ctx context.Context, msg string, fields ...zap.Field) {
	fields = TryAppendRequestIDFromContext(ctx, fields)
	l.l.Error(msg, fields...)
}
