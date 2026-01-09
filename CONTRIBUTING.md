# Contributing to BallBreaker ⚡

Thank you for your interest in improving **BallBreaker**! We follow the philosophy of "Stop the snowball," and your contributions help us keep our systems stable.

## 🛠️ Development Setup

1.  **Fork** the repository and clone it locally.
2.  Ensure you have **Go 1.25+** installed.
3.  We use a `Makefile` to simplify common tasks:

| Command | Description |
|---------|-------------|
| `make test` | Run all tests with `-race` detector |
| `make bench` | Run benchmarks to ensure no performance regressions |
| `make coverage` | Check test coverage |

## 🧪 Requirements for Pull Requests

To maintain the quality of the project, please ensure:

-   **Tests:** Every new feature or bug fix must include corresponding tests.
-   **Performance:** If you are touching the `Do` method or state transitions, run `make bench` to verify we still have **zero allocations**.
-   **Branching:** Create your feature branch off of the `dev` branch.
-   **Linter:** Run `golangci-lint` if possible, or at least ensure your code is formatted with `go fmt`.

## 🤝 Community Guidelines

-   **Be kind:** We follow the principle of "Be kind always."
-   **Open an Issue:** For large changes, please open an issue first to discuss the approach.

---
> **BallBreaker** — *Stop the snowball.*
