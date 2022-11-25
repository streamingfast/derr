package derr

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRetryContext(t *testing.T) {
	var count int
	err := RetryContext(context.Background(), 2, func(ctx context.Context) error {
		count++
		return fmt.Errorf("I failed")
	})
	assert.Error(t, err)
	assert.Equal(t, 3, count)
}

func TestRetryContextNextFailure(t *testing.T) {
	var count int
	err := RetryContext(context.Background(), 2, func(ctx context.Context) error {
		count++
		if count > 1 {
			return nil
		}
		return fmt.Errorf("I failed")
	})
	assert.NoError(t, err)
	assert.Equal(t, 2, count)
}

func TestRetryContextCanceledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := RetryContext(ctx, 2, func(ctx context.Context) error {
		t.Fail()
		return nil
	})
	assert.Error(t, err, "context cancelled")
}
