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

import "time"

// Option is a function that configures a CircuitBreaker.
type Option func(*CircuitBreaker)

// WithMaxFailures sets the number of failures before opening the circuit.
func WithMaxFailures(n int) Option {
	return func(cb *CircuitBreaker) {
		cb.maxFailures = n
	}
}

// WithMaxSuccesses sets the number of successes before closing the circuit.
func WithMaxSuccesses(n int) Option {
	return func(cb *CircuitBreaker) {
		cb.maxSuccesses = n
	}
}

// WithTimeout sets the time to wait before trying to close the circuit.
func WithTimeout(d time.Duration) Option {
	return func(cb *CircuitBreaker) {
		cb.timeout = d
	}
}

// WithOnStateChange sets the hook that is called when the circuit breaker changes state.
func WithOnStateChange(hook StateChangeHook) Option {
	return func(cb *CircuitBreaker) {
		cb.OnStateChange = hook
	}
}
