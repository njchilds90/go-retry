
// Package retry provides a retry policy for handling transient failures.
package retry

import (
    "context"
    "errors"
    "testing"
    "time"
)

// ErrMaxAttemptsReached is returned when the maximum number of attempts is reached.
var ErrMaxAttemptsReached = errors.New("max attempts reached")

// ErrContextCancelled is returned when the context is cancelled.
var ErrContextCancelled = errors.New("context cancelled")

// Policy is a retry policy for handling transient failures.
type Policy struct {
    maxAttempts int
    backoff     Backoff
}

// NewPolicy returns a new retry policy.
//
// Options:
//   - WithMaxAttempts: sets the maximum number of attempts.
//   - WithBackoff: sets the backoff strategy.
func NewPolicy(options ...PolicyOption) *Policy {
    policy := &Policy{
        maxAttempts: 3,
        backoff:     &ExponentialBackoff{},
    }

    for _, option := range options {
        option(policy)
    }

    return policy
}

// PolicyOption is an option for the NewPolicy function.
type PolicyOption func(*Policy)

// WithMaxAttempts sets the maximum number of attempts.
func WithMaxAttempts(maxAttempts int) PolicyOption {
    return func(policy *Policy) {
        policy.maxAttempts = maxAttempts
    }
}

// WithBackoff sets the backoff strategy.
func WithBackoff(backoff Backoff) PolicyOption {
    return func(policy *Policy) {
        policy.backoff = backoff
    }
}

// Backoff is a backoff strategy for retrying failed operations.
type Backoff interface {
    NextAttempt(attempt int) time.Duration
}

// ExponentialBackoff is an exponential backoff strategy.
type ExponentialBackoff struct {
    initial time.Duration
}

// NewExponentialBackoff returns a new exponential backoff strategy.
func NewExponentialBackoff(initial time.Duration) *ExponentialBackoff {
    return &ExponentialBackoff{initial: initial}
}

// NextAttempt returns the next attempt duration.
func (b *ExponentialBackoff) NextAttempt(attempt int) time.Duration {
    return b.initial * time.Duration(attempt*2)
}

// Retry retries the given operation according to the policy.
func (p *Policy) Retry(ctx context.Context, fn func(ctx context.Context) error) error {
    if fn == nil {
        return errors.New("fn is required")
    }

    for attempt := 1; attempt <= p.maxAttempts; attempt++ {
        if err := fn(ctx); err != nil {
            if errors.Is(err, context.Canceled) {
                return err
            }

            if attempt < p.maxAttempts {
                duration := p.backoff.NextAttempt(attempt)
                select {
                case <-ctx.Done():
                    return ctx.Err()
                case <-time.After(duration):
                }
            } else {
                return fmt.Errorf("max attempts reached: %w", err)
            }
        } else {
            return nil
        }
    }

    return ErrMaxAttemptsReached
}

func TestPolicy_Retry(t *testing.T) {
    tests := []struct {
        name       string
        policy     *Policy
        fn         func(ctx context.Context) error
        wantError error
    }{
        {
            name: "retry failed",
            policy: NewPolicy(
                WithMaxAttempts(3),
                WithBackoff(NewExponentialBackoff(500 * time.Millisecond)),
            ),
            fn: func(ctx context.Context) error {
                return errors.New("failed")
            },
            wantError: ErrMaxAttemptsReached,
        },
        {
            name: "retry succeeds",
            policy: NewPolicy(
                WithMaxAttempts(3),
                WithBackoff(NewExponentialBackoff(500 * time.Millisecond)),
            ),
            fn: func(ctx context.Context) error {
                return nil
            },
            wantError: nil,
        },
        {
            name: "retry context cancelled",
            policy: NewPolicy(
                WithMaxAttempts(3),
                WithBackoff(NewExponentialBackoff(500 * time.Millisecond)),
            ),
            fn: func(ctx context.Context) error {
                return errors.New("failed")
            },
            wantError: context.Canceled,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            ctx, cancel := context.WithCancel(context.Background())
            if tt.name == "retry context cancelled" {
                cancel()
            }

            err := tt.policy.Retry(ctx, tt.fn)

            if !errors.Is(err, tt.wantError) {
                t.Errorf("want error %v, got %v", tt.wantError, err)
            }
        })
    }
}

func BenchmarkPolicy_Retry(b *testing.B) {
    policy := NewPolicy(
        WithMaxAttempts(3),
        WithBackoff(NewExponentialBackoff(500 * time.Millisecond)),
    )

    b.ResetTimer()

    for i := 0; i < b.N; i++ {
        err := policy.Retry(context.Background(), func(ctx context.Context) error {
            return errors.New("failed")
        })

        if err != ErrMaxAttemptsReached {
            b.Errorf("want error %v, got %v", ErrMaxAttemptsReached, err)
        }
    }
}

func ExamplePolicy_Retry() {
    policy := NewPolicy(
        WithMaxAttempts(3),
        WithBackoff(NewExponentialBackoff(500 * time.Millisecond)),
    )

    err := policy.Retry(context.Background(), func(ctx context.Context) error {
        return errors.New("failed")
    })

    if err != nil {
        println(err)
    }
    // Output: max attempts reached: failed
    