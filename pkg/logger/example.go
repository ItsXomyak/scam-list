package logger

import (
	"context"
	"errors"
	"fmt"
)

// Example demonstrates how to use the logger with context and error wrapping
func Example() {
	ctx := context.Background()

	// Initialize logger
	l := InitLogger("test-service", LevelDebug)

	// Create a context with action and request ID
	ctx = WithAction(ctx, "example_action")
	// Add a request ID
	ctx = WithRequestID(ctx, "req-12345")

	// Log an info message
	l.Info(ctx, "This is an info message")

	// Simulate an error and wrap it with context
	if err := someLogic(ctx); err != nil {
		// Use ErrorCtx() to extract context from the error
		// and log the error with additional context
		// Lookup the log message, you can see both new action from someLogic function and request_id in the output
		l.Error(ErrorCtx(ctx, err), "failed to execute someLogic", err)
	}
}

// someLogic simulates a function that does some work and returns an error
func someLogic(ctx context.Context) error {
	// create a context with action field
	ctx = WithAction(ctx, "someLogic action")

	if err := secondLogic(ctx); err != nil {
		// Wrap the error with the current context
		return WrapError(ErrorCtx(ctx, err), fmt.Errorf("failed to execute secondLogic: %w", err))
	}

	return nil
}

func secondLogic(ctx context.Context) error {
	// create a context with action field
	ctx = WithAction(ctx, "secondLogic action")

	err := errors.New("secondLogic example error")
	// Wrap the error with the current context
	return WrapError(ctx, err)
}
