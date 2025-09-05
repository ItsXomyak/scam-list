package logger

import (
	"context"
	"errors"
)

// errorWithLogCtx is a custom error type that wraps an error and includes LogCtx
type errorWithLogCtx struct {
	err    error
	logCtx LogCtx
}

// errorWithLogCtx implements the error interface
func (e *errorWithLogCtx) Error() string {
	return e.err.Error()
}

// Unwrap allows unwrapping the original error
func (e *errorWithLogCtx) Unwrap() error {
	return e.err
}

// WrapError wraps an error with the current LogCtx from the context
func WrapError(ctx context.Context, err error) error {
	if err == nil {
		return nil
	}

	// If already wrapped, just update logCtx
	var e *errorWithLogCtx
	if errors.As(err, &e) {
		if x, ok := ctx.Value(logCtxKey).(LogCtx); ok {
			e.logCtx = x
		}
		e.err = err
		return e
	}

	c := LogCtx{}
	if x, ok := ctx.Value(logCtxKey).(LogCtx); ok {
		c = x
	}
	return &errorWithLogCtx{
		err:    err,
		logCtx: c,
	}
}

// ErrorCtx extracts the LogCtx from an error if it is of type errorWithLogCtx
func ErrorCtx(ctx context.Context, err error) context.Context {
	var e *errorWithLogCtx
	if errors.As(err, &e) && e != nil {
		return context.WithValue(ctx, logCtxKey, e.logCtx)
	}
	return ctx
}
