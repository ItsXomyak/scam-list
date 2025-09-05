package logger

import (
	"context"
	"log/slog"
	"os"
	"time"
)

const (
	LevelDebug string = "DEBUG"
	LevelInfo  string = "INFO"
	LevelWarn  string = "WARN"
	LevelError string = "ERROR"
)

type Logger interface {
	Debug(ctx context.Context, msg string, args ...any)
	Info(ctx context.Context, msg string, args ...any)
	Warn(ctx context.Context, msg string, args ...any)
	Error(ctx context.Context, msg string, err error, args ...any)
	GetSlogLogger() *slog.Logger
}

type logger struct {
	slog *slog.Logger
}

// Initialize logger with service name and log level
func InitLogger(serviceName, logLevel string) Logger {
	level := new(slog.LevelVar)
	switch logLevel {
	case LevelDebug:
		level.Set(slog.LevelDebug)
	case LevelInfo:
		level.Set(slog.LevelInfo)
	case LevelWarn:
		level.Set(slog.LevelWarn)
	case LevelError:
		level.Set(slog.LevelError)
	default:
		level.Set(slog.LevelDebug)
	}

	// Custom handler
	handler := &contextHandler{
		handler: slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: level,
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				// Format time as ISO 8601
				if a.Key == slog.TimeKey {
					if t, ok := a.Value.Any().(time.Time); ok {
						return slog.Attr{Key: "timestamp", Value: slog.StringValue(t.Format(time.RFC3339))}
					}
				}
				return a
			},
			AddSource: false, // add a SourceKey to attrs with the location of the log call file
		}),
	}

	// Create base logger with service and hostname
	base := slog.New(handler).With(
		slog.String("service", serviceName),
	)

	return &logger{
		slog: base,
	}
}

// Context handler to inject values from context
type contextHandler struct {
	handler slog.Handler
}

func (h *contextHandler) Enabled(ctx context.Context, lvl slog.Level) bool {
	return h.handler.Enabled(ctx, lvl)
}

func (h *contextHandler) Handle(ctx context.Context, r slog.Record) error {
	if c, ok := ctx.Value(logCtxKey).(LogCtx); ok {
		if c.Action != "" {
			r.AddAttrs(slog.String("action", c.Action))
		}
		if c.RequestID != "" {
			r.AddAttrs(slog.String("request_id", c.RequestID))
		}
	}

	return h.handler.Handle(ctx, r)
}

func (h *contextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &contextHandler{handler: h.handler.WithAttrs(attrs)}
}

func (h *contextHandler) WithGroup(name string) slog.Handler {
	return &contextHandler{handler: h.handler.WithGroup(name)}
}

// Logger methods
func (l *logger) Debug(ctx context.Context, msg string, args ...any) {
	l.slog.DebugContext(ctx, msg, args...)
}

func (l *logger) Info(ctx context.Context, msg string, args ...any) {
	l.slog.InfoContext(ctx, msg, args...)
}

func (l *logger) Warn(ctx context.Context, msg string, args ...any) {
	l.slog.WarnContext(ctx, msg, args...)
}

func (l *logger) Error(ctx context.Context, msg string, err error, args ...any) {
	attrs := []any{
		"error", err.Error(),
	}
	attrs = append(attrs, args...)
	l.slog.ErrorContext(ctx, msg, attrs...)
}

func (l *logger) GetSlogLogger() *slog.Logger {
	return l.slog
}

// ValidateLogLevel validates if the given string is valid logger level(DEBUG, INFO, WARN, ERROR).
func ValidateLogLevel(lvl string) bool {
	switch lvl {
	case LevelDebug, LevelError, LevelWarn, LevelInfo:
		return true
	default:
		return false
	}
}
