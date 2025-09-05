package logger

import (
	"context"
)

type (
	// LogCtx holds contextual information for logging
	LogCtx struct {
		Action    string
		RequestID string
	}

	// logCtxKeyStruct is an unexported type for context keys defined in this package.
	logCtxKeyStruct struct{}
)

// logCtxKey is the key for log context values
var logCtxKey = &logCtxKeyStruct{}

// WithLogCtx returns a new context with the provided LogCtx
func WithLogCtx(ctx context.Context, newLc LogCtx) context.Context {
	// Check if there's an existing LogCtx and merge values
	if lc, ok := ctx.Value(logCtxKey).(LogCtx); ok {
		if newLc.Action == "" {
			newLc.Action = lc.Action
		}
		if newLc.RequestID == "" {
			newLc.RequestID = lc.RequestID
		}

		return context.WithValue(ctx, logCtxKey, newLc)
	}
	return context.WithValue(ctx, logCtxKey, newLc)
}

// WithRequestID adds or updates the RequestID in the LogCtx within the context
func WithRequestID(ctx context.Context, requestID string) context.Context {
	if lc, ok := ctx.Value(logCtxKey).(LogCtx); ok {
		lc.RequestID = requestID
		return context.WithValue(ctx, logCtxKey, lc)
	}
	return context.WithValue(ctx, logCtxKey, LogCtx{RequestID: requestID})
}

// WithAction adds or updates the Action in the LogCtx within the context
func WithAction(ctx context.Context, action string) context.Context {
	if lc, ok := ctx.Value(logCtxKey).(LogCtx); ok {
		lc.Action = action
		return context.WithValue(ctx, logCtxKey, lc)
	}
	return context.WithValue(ctx, logCtxKey, LogCtx{Action: action})
}
