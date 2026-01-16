/*
Copyright 2026 Diego Cena <diego.cena@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package ballbreaker

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

type CircuitStateType int

// StateChangeHook is a callback function that is called when the circuit breaker changes state.
type StateChangeHook func(from, to CircuitStateType)

const (
	StateClosed CircuitStateType = iota
	StateOpen
	StateHalfOpen
)

var ErrCircuitBreakerOpen = errors.New("circuit breaker is open")

type CircuitBreaker struct {
	sync.Mutex

	state          CircuitStateType
	maxFailures    int
	failureCount   int
	successCount   int
	maxSuccesses   int
	timeout        time.Duration
	lastFailure    time.Time
	openStateError error

	// OnStateChange is called whenever the circuit breaker transitions between states.
	OnStateChange StateChangeHook
}

// New creates a new circuit breaker with functional options.
// If no options are provided, it uses sensible defaults:
// - MaxFailures: 5
// - MaxSuccesses: 2
// - Timeout: 5 seconds
func New(opts ...Option) *CircuitBreaker {
	cb := &CircuitBreaker{
		state:        StateClosed,
		maxFailures:  5,
		maxSuccesses: 2,
		timeout:      5 * time.Second,
		lastFailure:  time.Now(),
	}

	for _, opt := range opts {
		opt(cb)
	}

	return cb
}

// State returns the current state of the circuit breaker.
func (cb *CircuitBreaker) State() CircuitStateType {
	cb.Lock()
	defer cb.Unlock()
	if cb.state == StateOpen && time.Since(cb.lastFailure) > cb.timeout {
		return StateHalfOpen
	}
	return cb.state
}

// Do executes the given function and returns the result.
// It is a shorthand for DoWithContext(context.Background(), fn).
func (cb *CircuitBreaker) Do(fn func() error) error {
	return cb.DoWithContext(context.Background(), fn)
}

// DoWithContext executes the given function and returns the result.
// If the context is cancelled, it returns the context error.
// If the circuit breaker is open, it returns the last recorded error.
func (cb *CircuitBreaker) DoWithContext(ctx context.Context, fn func() error) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	cb.Lock()

	switch cb.state {
	case StateClosed:
		cb.Unlock()
		err := fn()
		if ctx.Err() != nil {
			return ctx.Err()
		}
		cb.Lock()
		defer cb.Unlock()
		return cb.handleClosedState(err)

	case StateOpen:
		if time.Since(cb.lastFailure) > cb.timeout {
			cb.setState(StateHalfOpen)
			cb.Unlock()
			err := fn()
			if ctx.Err() != nil {
				return ctx.Err()
			}
			cb.Lock()
			defer cb.Unlock()
			return cb.handleHalfOpenState(err)
		}
		defer cb.Unlock()
		//return cb.openStateError
		return fmt.Errorf("%w: %w", ErrCircuitBreakerOpen, cb.openStateError)

	case StateHalfOpen:
		cb.Unlock()
		err := fn()
		if ctx.Err() != nil {
			return ctx.Err()
		}
		cb.Lock()
		defer cb.Unlock()
		return cb.handleHalfOpenState(err)
	}

	cb.Unlock()
	return nil
}

// handleClosedState handles the closed state of the circuit breaker.
// Must be called with lock held.
func (cb *CircuitBreaker) handleClosedState(err error) error {
	if err != nil {
		cb.failureCount++
		if cb.failureCount >= cb.maxFailures {
			cb.setState(StateOpen)
			cb.openStateError = err
			cb.failureCount = 0
			cb.lastFailure = time.Now()
		}
		return err
	}
	cb.failureCount = 0
	return nil
}

// handleHalfOpenState handles the half open state of the circuit breaker.
// Must be called with lock held.
func (cb *CircuitBreaker) handleHalfOpenState(err error) error {
	if err != nil {
		cb.setState(StateOpen)
		cb.openStateError = err
		cb.failureCount = 0
		cb.successCount = 0
		cb.lastFailure = time.Now()
		return err
	}

	cb.successCount++
	if cb.successCount >= cb.maxSuccesses {
		cb.setState(StateClosed)
		cb.successCount = 0
		cb.failureCount = 0
		cb.openStateError = nil
	}
	return nil
}

// setState changes the state of the circuit breaker and triggers the OnStateChange hook.
// It assumes the lock is held but may temporarily release it to call the hook.
func (cb *CircuitBreaker) setState(newState CircuitStateType) {
	if cb.state == newState {
		return
	}

	oldState := cb.state
	cb.state = newState

	if cb.OnStateChange != nil {
		hook := cb.OnStateChange
		// We call the hook while holding the lock to ensure state consistency,
		// but we could also unlock/lock. However, in this library's simple usage,
		// calling it inside is standard unless we expect heavy/complex hooks.
		// For safety against deadlocks in user code, we recommend simple hooks.
		hook(oldState, newState)
	}
}
