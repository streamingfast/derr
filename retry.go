package derr

import (
	"context"
	"errors"
	"fmt"
	"time"

	retry "github.com/sethvargo/go-retry"
)

type FatalError struct {
	original error
}

// NewFatalError creates a new [FatalError] struct ensuring `original` error is non-nil
// otherwise this function panics with an error.
func NewFatalError(original error) *FatalError {
	if original == nil {
		panic(fmt.Errorf("the 'original' argument is mandatory"))
	}

	return &FatalError{original}
}

func (r *FatalError) Unwrap() error {
	return r.original
}

func (r *FatalError) Error() string {
	return r.original.Error()
}

// RetryableError can be returned by your handler either [SinkerHandlers#HandleBlockScopedData] or
// [SinkerHandlers#HandleBlockUndoSignal] to notify the sinker that it's a retryable error and the
// stream can continue
type RetryableError struct {
	original error
}

// NewRetryableError creates a new [RetryableError] struct ensuring `original` error is non-nil
// otherwise this function panics with an error.
func NewRetryableError(original error) *RetryableError {
	if original == nil {
		panic(fmt.Errorf("the 'original' argument is mandatory"))
	}

	return &RetryableError{original}
}

func (r *RetryableError) Unwrap() error {
	return r.original
}

func (r *RetryableError) Error() string {
	return fmt.Sprintf("%s (retryable)", r.original)
}

// Retry re-executes the function `f` if it returns an error. If you return a  `derr.FatalError` your function
// will not be retried.
func Retry(retries uint64, f func(ctx context.Context) error) error {
	return RetryContext(context.Background(), retries, f)
}

// RetryContext re-executes the function `f` if it returns an error. If you return a  `derr.FatalError` your function
// will not be retried.
func RetryContext(ctx context.Context, retries uint64, f func(ctx context.Context) error) error {
	return retry.Do(ctx, backoff(retries), func(ctx context.Context) error {
		err := f(ctx)
		if err != nil {
			var fatalError *FatalError
			if errors.As(err, &fatalError) {
				return fatalError.original
			}
			return retry.RetryableError(err)
		}
		return nil
	})
}

func backoff(maxretries uint64) retry.Backoff {
	b := retry.NewFibonacci(time.Second)
	b = retry.WithMaxRetries(maxretries, b)
	b = retry.WithCappedDuration(5*time.Second, b)
	return b
}
