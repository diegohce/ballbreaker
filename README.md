# ⚡ BallBreaker
> *Stop the snowball.*

[![Go Reference](https://pkg.go.dev/badge/github.com/diegohce/ballbreaker.svg)](https://pkg.go.dev/github.com/diegohce/ballbreaker)
[![Go Report Card](https://goreportcard.com/badge/github.com/diegohce/ballbreaker)](https://goreportcard.com/report/github.com/diegohce/ballbreaker)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

**BallBreaker** is a high-performance, thread-safe Circuit Breaker pattern implementation for Go. It helps you prevent cascading failures in distributed systems by providing a robust mechanism to gracefully handle service degradation.

## ✨ Features

- 🔒 **Thread-Safe**: Built with concurrency in mind using `sync.Mutex`.
- 🚀 **Performant**: Zero-allocation state transitions during steady state.
- 🛠️ **Configurable**: Thresholds for failures, successes, and recovery timeouts.
- 📊 **State Inspection**: Real-time monitoring of the circuit status.
- ✅ **Tested**: ~95% code coverage and verified with `-race` detector.

## 🕹️ Quick Start

### Installation

```bash
go get github.com/diegohce/ballbreaker
```

### Usage

```go
package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/diegohce/ballbreaker"
)

func main() {
	// Create a new breaker: 
	// - 3 failures to open
	// - 2 successes to close
	// - 5 seconds timeout before trying to recover
	cb := ballbreaker.New(3, 2, 5*time.Second)

	err := cb.Do(func() error {
		// Your potentially failing logic here
		resp, err := http.Get("http://example.com")
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		
		if resp.StatusCode >= 500 {
			return fmt.Errorf("server error: %d", resp.StatusCode)
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Operation failed or circuit is open: %v\n", err)
	}
}
```

## 📐 How it works

The circuit breaker has three states:

1.  **Closed**: The normal state. Requests flow normally. If failures exceed the `maxFailures` threshold, the circuit **Opens**.
2.  **Open**: Requests fail immediately without executing the function. After the `timeout` expires, the next request will transition the circuit to **Half-Open**.
3.  **Half-Open**: A limited number of trial requests are allowed. If `maxSuccesses` are reached, the circuit **Closes**. If any request fails, it reverts to **Open**.

## 🍕 Inspiration

The name **BallBreaker** is a tribute to the fictional arcade game featured in the restaurant  "The Original Beef of Chicagoland" from the TV series ***The Bear***. According to Matty Matheson's character Fak, Ballbreaker is a "Norwegian knock-off of Mortal Kombat."

## 🛡️ License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

---
> **BallBreaker** — *Stop the snowball.*
