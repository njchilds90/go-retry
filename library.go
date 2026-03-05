package retry

import (
    "context"
    "errors"
    "time"
)

// Policy represents a retry policy.
type Policy struct {
    maxAttempts int
    backoff     Backoff
}

// NewPolicy creates a new retry policy.
func NewPolicy(opts ...Option) *Policy {
    policy := &Policy{
        maxAttempts: 3,
        backoff:     ExponentialBackoff(500 * time.Millisecond),
    }

    for _, opt := range opts {
        opt(policy)
    }

    return policy
}

// Retry applies the retry policy to a function.
func (p *Policy) Retry(ctx context.Context, fn func(ctx context.Context) error) error {
    var attempt int
    var lastErr error

    for attempt < p.maxAttempts {
        err := fn(ctx)
        if err == nil {
            return nil
        }

        lastErr = err
        attempt++

        if attempt < p.maxAttempts {
            delay := p.backoff.CalculateDelay(attempt)
            select {
            case <-ctx.Done():
                return ctx.Err()
            case <-time.After(delay):
            }
        }
    }

    return lastErr
}

// Option represents a retry policy option.
type Option func(*Policy)

// WithMaxAttempts sets the maximum number of attempts.
func WithMaxAttempts(maxAttempts int) Option {
    return func(p *Policy) {
        p.maxAttempts = maxAttempts
    }
}

// WithBackoff sets the backoff strategy.
func WithBackoff(backoff Backoff) Option {
    return func(p *Policy) {
        p.backoff = backoff
    }
}

// Backoff represents a backoff strategy.
type Backoff interface {
    CalculateDelay(attempt int) time.Duration
}

// ExponentialBackoff creates an exponential backoff strategy.
func ExponentialBackoff(initialDuration time.Duration) Backoff {
    return &exponentialBackoff{
        initialDuration: initialDuration,
    }
}

type exponentialBackoff struct {
    initialDuration time.Duration
}

func (e *exponentialBackoff) CalculateDelay(attempt int) time.Duration {
    return e.initialDuration * (2 << uint(attempt-1))
}

// LinearBackoff creates a linear backoff strategy.
func LinearBackoff(initialDuration time.Duration) Backoff {
    return &linearBackoff{
        initialDuration: initialDuration,
    }
}

type linearBackoff struct {
    initialDuration time.Duration
}

func (l *linearBackoff) CalculateDelay(attempt int) time.Duration {
    return l.initialDuration * time.Duration(attempt)
}
