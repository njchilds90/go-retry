package retry

// Package retry provides a flexible retry mechanism with various backoff strategies.
//
// See https://pkg.go.dev/github.com/njchilds90/go-retry for documentation and examples.
import (
	"context"
	"errors"
	"time"
)

// ErrMaxAttempts is returned when the maximum number of attempts is exceeded.
var ErrMaxAttempts = errors.New("maximum number of attempts exceeded")

// ErrInvalidOption is returned when an invalid option is provided.
var ErrInvalidOption = errors.New("invalid option provided")

// Policy represents a retry policy.
//
// Policy defines the maximum number of attempts and the backoff strategy to use when retrying a function.
type Policy struct {
	maxAttempts int
	backoff     Backoff
}

// NewPolicy creates a new retry policy with the given options.
//
// NewPolicy returns an error if any of the options are invalid.
func NewPolicy(opts ...Option) (*Policy, error) {
	if len(opts) == 0 {
		return nil, ErrInvalidOption
	}

	policy := &Policy{
		maxAttempts: 3,
		backoff:     ExponentialBackoff(500 * time.Millisecond),
	}

	for _, opt := range opts {
		opt(policy)
	}

	if policy.maxAttempts < 1 {
		return nil, ErrInvalidOption
	}

	return policy, nil
}

// Retry applies the retry policy to a function.
//
// Retry calls the given function with the provided context and retries it according to the policy.
//
// If the function returns an error, Retry will wait for the calculated delay before retrying.
//
// If the maximum number of attempts is exceeded, Retry returns ErrMaxAttempts.
func (p *Policy) Retry(ctx context.Context, fn func(ctx context.Context) error) error {
	if ctx == nil {
		return errors.New("context is nil")
	}

	if fn == nil {
		return errors.New("function is nil")
	}

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

	return fmt.Errorf("retry failed after %d attempts: %w", p.maxAttempts, lastErr)
}

// Option represents a retry policy option.
//
// Option is a functional option that can be used to customize the retry policy.
type Option func(*Policy)

// WithMaxAttempts sets the maximum number of attempts.
//
// WithMaxAttempts returns an error if the given number of attempts is less than 1.
func WithMaxAttempts(maxAttempts int) Option {
	return func(p *Policy) {
		p.maxAttempts = maxAttempts
	}
}

// WithBackoff sets the backoff strategy.
//
// WithBackoff returns an error if the given backoff strategy is nil.
func WithBackoff(backoff Backoff) Option {
	return func(p *Policy) {
		p.backoff = backoff
	}
}

// Backoff represents a backoff strategy.
//
// Backoff defines the delay between attempts.
type Backoff interface {
	CalculateDelay(attempt int) time.Duration
}

// ExponentialBackoff creates an exponential backoff strategy.
//
// ExponentialBackoff returns a backoff strategy that doubles the delay after each attempt.
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
