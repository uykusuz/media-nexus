package util

import (
	"context"
	"media-nexus/logger"
)

type ContextValue string

const (
	contextLogger ContextValue = "context"
)

func WithLogger(ctx context.Context, logger logger.Logger) context.Context {
	return context.WithValue(ctx, contextLogger, logger)
}

func Logger(ctx context.Context) logger.Logger {
	return ctx.Value(contextLogger).(logger.Logger)
}
