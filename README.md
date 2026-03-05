# Go-Retry

A simple and efficient retry library for Go applications.

## Overview

Go-retry provides a simple way to add retry mechanisms to your Go code. It allows you to define a retry policy and apply it to any function that returns an error.

## Installation

To install go-retry, run the following command:

```go
import "github.com/go-retry/go-retry"
```

## Usage

Here is an example of how to use go-retry:

    package main

    import (
        "context"
        "fmt"
        "time"

        "github.com/go-retry/go-retry"
    )

    func main() {
        // Define a retry policy
        policy := retry.NewPolicy(
            retry.WithMaxAttempts(3),
            retry.WithBackoff(retry.ExponentialBackoff(500 * time.Millisecond)),
        )

        // Apply the retry policy to a function
        err := policy.Retry(context.Background(), func(ctx context.Context) error {
            // Simulate a failing function
            return fmt.Errorf("failed")
        })

        if err != nil {
            fmt.Println(err)
        }
    }

## API Reference

### Policy

* `NewPolicy(opts ...Option) *Policy`: Creates a new retry policy.
* `Retry(ctx context.Context, fn func(ctx context.Context) error) error`: Applies the retry policy to a function.

### Options

* `WithMaxAttempts(maxAttempts int) Option`: Sets the maximum number of attempts.
* `WithBackoff(backoff Backoff) Option`: Sets the backoff strategy.

### Backoff

* `ExponentialBackoff(initialDuration time.Duration) Backoff`: Creates an exponential backoff strategy.
* `LinearBackoff(initialDuration time.Duration) Backoff`: Creates a linear backoff strategy.
