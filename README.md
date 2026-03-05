# Go-Retry

A simple and efficient retry library for Go applications, utilizing Go's idioms and best practices, including functional options, interface-first design, and context-aware I/O operations.

## Overview

Go-retry provides a simple way to add retry mechanisms to your Go code. It allows you to define a retry policy and apply it to any function that returns an error.

## Installation

To install go-retry, run the following command:
```bash
go get github.com/njchilds90/go-retry
```

[![PkgGoDev](https://pkg.go.dev/badge/github.com/njchilds90/go-retry)](https://pkg.go.dev/github.com/njchilds90/go-retry)

## Usage

Here is an example of how to use go-retry:
```go
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/njchilds90/go-retry"
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
```

## API Reference

### Policy

* `NewPolicy(opts ...Option) *Policy`: Creates a new retry policy. The policy is used to define the retry strategy, including the maximum number of attempts and the backoff strategy.
* `Retry(ctx context.Context, fn func(ctx context.Context) error) error`: Applies the retry policy to a function. This method will attempt to execute the provided function, retrying as necessary according to the defined policy.

### Options

* `WithMaxAttempts(maxAttempts int) Option`: Sets the maximum number of attempts. This option is used to define the maximum number of times the retry policy will attempt to execute the provided function.
* `WithBackoff(backoff Backoff) Option`: Sets the backoff strategy. This option is used to define the backoff strategy, which determines how long the retry policy will wait between attempts.

### Backoff

* `ExponentialBackoff(initialDuration time.Duration) Backoff`: Creates an exponential backoff strategy. This backoff strategy will increase the wait time exponentially between attempts.
* `LinearBackoff(initialDuration time.Duration) Backoff`: Creates a linear backoff strategy. This backoff strategy will increase the wait time linearly between attempts.