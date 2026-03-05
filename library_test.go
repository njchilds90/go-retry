package retry

import (
    "context"
    "testing"
    "time"
)

func TestPolicy_Retry(t *testing.T) {
    policy := NewPolicy(
        WithMaxAttempts(3),
        WithBackoff(ExponentialBackoff(500 * time.Millisecond)),
    )

    err := policy.Retry(context.Background(), func(ctx context.Context) error {
        return errors.New("failed")
    })

    if err == nil {
        t.Errorf("expected error, got nil")
    }
}

func TestPolicy_Retry_Succeeds(t *testing.T) {
    policy := NewPolicy(
        WithMaxAttempts(3),
        WithBackoff(ExponentialBackoff(500 * time.Millisecond)),
    )

    err := policy.Retry(context.Background(), func(ctx context.Context) error {
        return nil
    })

    if err != nil {
        t.Errorf("expected nil, got error: %v", err)
    }
}

func TestPolicy_Retry_ContextCancel(t *testing.T) {
    policy := NewPolicy(
        WithMaxAttempts(3),
        WithBackoff(ExponentialBackoff(500 * time.Millisecond)),
    )

    ctx, cancel := context.WithCancel(context.Background())
    cancel()

    err := policy.Retry(ctx, func(ctx context.Context) error {
        return errors.New("failed")
    })

    if err != context.Canceled {
        t.Errorf("expected context.Canceled, got %v", err)
    }
}
