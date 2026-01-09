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
	"sync"
	"time"
)

type CircuitStateType int

const (
	StateClosed CircuitStateType = iota
	StateOpen
	StateHalfOpen
)

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
}

// New creates a new circuit breaker.
// maxFailures is the number of failures before opening the circuit.
// maxSuccesses is the number of successes before closing the circuit.
// timeout is the time to wait before trying to close the circuit.
func New(maxFailures int, maxSuccesses int, timeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		state:        StateClosed,
		maxFailures:  maxFailures,
		failureCount: 0,
		successCount: 0,
		maxSuccesses: maxSuccesses,
		timeout:      timeout,
		lastFailure:  time.Now(),
	}
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
// If the circuit breaker is open, it returns the last recorded error.
func (cb *CircuitBreaker) Do(fn func() error) error {
	cb.Lock()

	switch cb.state {
	case StateClosed:
		cb.Unlock()
		err := fn()
		cb.Lock()
		defer cb.Unlock()
		return cb.handleClosedState(err)

	case StateOpen:
		if time.Since(cb.lastFailure) > cb.timeout {
			cb.state = StateHalfOpen
			cb.Unlock()
			err := fn()
			cb.Lock()
			defer cb.Unlock()
			return cb.handleHalfOpenState(err)
		}
		defer cb.Unlock()
		return cb.openStateError

	case StateHalfOpen:
		cb.Unlock()
		err := fn()
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
			cb.state = StateOpen
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
		cb.state = StateOpen
		cb.openStateError = err
		cb.failureCount = 0
		cb.successCount = 0
		cb.lastFailure = time.Now()
		return err
	}

	cb.successCount++
	if cb.successCount >= cb.maxSuccesses {
		cb.state = StateClosed
		cb.successCount = 0
		cb.failureCount = 0
		cb.openStateError = nil
	}
	return nil
}
