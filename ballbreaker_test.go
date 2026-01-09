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
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

func TestBreaker(t *testing.T) {

	ErrUnexpectedStatusCode := errors.New("unexpected status code")
	serverStatus := http.StatusOK

	remote := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(serverStatus)
		fmt.Fprint(w, "OK")
	}))
	defer remote.Close()

	cb := New(2, 2, 1*time.Second)

	cases := []struct {
		name          string
		status        int
		expectedError error
		wait          time.Duration
		expectedState CircuitStateType
	}{
		{"initial_success", http.StatusOK, nil, 0, StateClosed},
		{"first_failure", http.StatusInternalServerError, ErrUnexpectedStatusCode, 0, StateClosed},
		{"second_failure_opens", http.StatusInternalServerError, ErrUnexpectedStatusCode, 0, StateOpen},
		{"fast_fail_while_open", http.StatusOK, ErrUnexpectedStatusCode, 0, StateOpen},
		{"recover_to_half_open_then_fail", http.StatusInternalServerError, ErrUnexpectedStatusCode, 1100 * time.Millisecond, StateOpen},
		{"still_open_after_half_open_failure", http.StatusOK, ErrUnexpectedStatusCode, 0, StateOpen},
		{"recover_to_half_open_success_1", http.StatusOK, nil, 1100 * time.Millisecond, StateHalfOpen},
		{"half_open_success_2_closes", http.StatusOK, nil, 0, StateClosed},
		{"final_success", http.StatusOK, nil, 0, StateClosed},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			serverStatus = c.status
			time.Sleep(c.wait)

			cbErr := cb.Do(func() error {
				res, err := http.Get(remote.URL)
				if err != nil {
					return err
				}
				defer res.Body.Close()
				if res.StatusCode != http.StatusOK {
					return ErrUnexpectedStatusCode
				}
				return nil
			})

			if cbErr != c.expectedError {
				t.Errorf("Expected error %v, got %v", c.expectedError, cbErr)
			}

			if state := cb.State(); state != c.expectedState {
				t.Errorf("Expected state %v, got %v", c.expectedState, state)
			}
		})
	}
}

func TestBreakerConcurrency(t *testing.T) {
	const workers = 20
	const iterations = 100

	cb := New(10, 10, 1*time.Second)

	var wg sync.WaitGroup
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				_ = cb.Do(func() error {
					time.Sleep(1 * time.Millisecond)
					return nil
				})
			}
		}()
	}
	wg.Wait()
}
