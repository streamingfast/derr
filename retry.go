package derr

import (
	"context"
	"errors"
	"time"

	retry "github.com/sethvargo/go-retry"
)

func backoff(maxretries uint64) retry.Backoff {
	b := retry.NewFibonacci(time.Second)
	b = retry.WithMaxRetries(maxretries, b)
	b = retry.WithCappedDuration(5*time.Second, b)
	return b
}

func Retry(retries uint64, f func(ctx context.Context) error) error {
	return RetryContext(context.Background(), retries, f)
}

func RetryContext(ctx context.Context, retries uint64, f func(ctx context.Context) error) error {
	err := retry.Do(ctx, backoff(retries), func(ctx context.Context) error {
		err := f(ctx)
		return retry.RetryableError(err)
	})
	return errors.Unwrap(err)
}
