# 🗺️ BallBreaker Roadmap

This document outlines the planned features and future direction for **BallBreaker**. We want to evolve from a lean circuit breaker into a comprehensive resilience toolkit for Go distributed systems.

## 🟢 Phase 1: Core Enhancements (Short-term)

- [x] **Context Support**: Implement `DoWithContext(ctx, fn)` to allow cancellation during execution.
- [x] **Functional Options**: Refactor the constructor to use the functional options pattern for cleaner configuration.
- [ ] **Customizable Failure Criteria**: Allow users to define what counts as a failure (e.g., specific HTTP status codes or error types).

## 🟡 Phase 2: Advanced Resilience (Mid-term)

- [ ] **Exponential Backoff**: Instead of a fixed timeout, implement backoff strategies for the transition from Open to Half-Open.
- [ ] **Sliding Window Counters**: Replace simple counters with a sliding window (time-based or count-based) for more accurate failure detection.
- [ ] **Concurrency Limiters**: Integrate bulkhead patterns to limit the number of concurrent requests to a specific service.

## 🔵 Phase 3: Observability & Ecosystem (Long-term)

- [ ] **Metrics Exporters**:
    - [ ] Native Prometheus collector.
    - [ ] OpenTelemetry integration for tracing and metrics.
- [ ] **Middleware Library**: Ready-to-use middleware for:
    - [ ] `net/http` (Standard library)
    - [ ] gRPC Interceptors
- [ ] **Dashboard Support**: A simple sidecar or web UI to visualize the state of multiple breakers in real-time.

---

### 💡 Have an idea?
If you think something is missing, feel free to [open an issue](https://github.com/diegohce/ballbreaker/issues) or discuss it in our [Contributing guide](CONTRIBUTING.md).

---
> **BallBreaker** — *Stop the snowball.*
