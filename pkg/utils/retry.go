package utils

import (
	"context"
	"time"

	"sentinel/pkg/logger"
)

const (
	// DefaultMaxRetries is the default maximum number of retry attempts
	DefaultMaxRetries = 3
	// DefaultRetryDelay is the default delay between retry attempts
	DefaultRetryDelay = 2 * time.Second
	// DefaultTimeout is the default timeout for operations
	DefaultTimeout = 10 * time.Second
)

// RetryConfig holds configuration for retry operations
type RetryConfig struct {
	MaxRetries int
	RetryDelay time.Duration
	Timeout    time.Duration
}

// DefaultRetryConfig returns a default retry configuration
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxRetries: DefaultMaxRetries,
		RetryDelay: DefaultRetryDelay,
		Timeout:    DefaultTimeout,
	}
}

// RetryFunc is a function that can be retried
type RetryFunc func(ctx context.Context) error

// RetryWithContext executes a function with retry logic and context timeout
func RetryWithContext(ctx context.Context, config RetryConfig, fn RetryFunc, operationName string) error {
	var lastErr error

	for attempt := 1; attempt <= config.MaxRetries; attempt++ {
		// Create context with timeout for this attempt
		attemptCtx, cancel := context.WithTimeout(ctx, config.Timeout)

		// Execute the function
		err := fn(attemptCtx)
		cancel()

		if err == nil {
			return nil
		}

		lastErr = err
		logger.CLogger.Errorf("ERROR: %s - Attempt %d/%d failed: %v", operationName, attempt, config.MaxRetries, err)

		// Don't sleep after the last attempt
		if attempt < config.MaxRetries {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(config.RetryDelay):
			}
		}
	}

	return lastErr
}
