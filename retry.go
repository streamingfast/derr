package derr

import (
	"context"
	"errors"
	"time"

	"github.com/sethvargo/go-retry"
)

func backoff(maxretries uint64) retry.Backoff {
	b, err := retry.NewFibonacci(time.Second)
	if err != nil {
		panic(err)
	}
	b = retry.WithMaxRetries(maxretries, b)
	b = retry.WithCappedDuration(5*time.Second, b)
	return b
}

func Retry(retries uint64, f func(ctx context.Context) error) error {
	return RetryContext(context.Background(), retries, f)
}

func RetryContext(ctx context.Context, retries uint64, f func(ctx context.Context) error) error {
	err := retry.Do(context.Background(), backoff(retries), func(ctx context.Context) error {
		err := f(ctx)
		return retry.RetryableError(err)
	})
	return errors.Unwrap(err)
}
