package kzklogger

import (
	"context"

	loggerconstant "github.com/kazekim/kzklogger/constant"
)

func getFromContext(ctx context.Context, key string) (value string, ok bool) {
	value, ok = ctx.Value(key).(string)
	return
}

func (g Logger) ApplyContext(ctx context.Context) *Logger {
	log := &g
	if value, ok := getFromContext(ctx, loggerconstant.ContextRequestId); ok {
		log = log.WithRequestID(value)
	}

	if value, ok := getFromContext(ctx, loggerconstant.ContextUserId); ok {
		log = log.WithUserID(value)
	}
	return log
}
