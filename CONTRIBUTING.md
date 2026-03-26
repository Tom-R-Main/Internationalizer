# Contributing to Internationalizer

Thank you for your interest in contributing! This document explains how to get started.

## Developer Certificate of Origin (DCO)

All contributions must be signed off, certifying that you wrote the code or have the right to submit it under the AGPL-3.0 license. Add a sign-off line to your commits:

```
git commit -s -m "Your commit message"
```

This adds a `Signed-off-by: Your Name <your@email.com>` line to your commit message.

## Development Setup

1. Install Go 1.22 or later
2. Clone the repository:
   ```
   git clone https://github.com/Tom-R-Main/Internationalizer.git
   cd Internationalizer
   ```
3. Build:
   ```
   go build ./cmd/internationalizer
   ```
4. Run tests:
   ```
   go test ./... -race
   ```

## Code Style

- Run `gofmt` on all Go files before committing
- Run `go vet ./...` to catch common issues
- Follow standard Go conventions and idioms
- Keep error messages lowercase (Go convention)
- Prefer returning errors over panicking

## Pull Request Process

1. Fork the repository and create a feature branch
2. Make your changes with clear, descriptive commit messages
3. Ensure all tests pass: `go test ./... -race`
4. Ensure code compiles cleanly: `go build ./...`
5. Sign off all commits (DCO requirement)
6. Open a pull request with a clear description of what and why

## Reporting Issues

Open an issue on GitHub with:
- What you expected to happen
- What actually happened
- Steps to reproduce
- Your Go version and OS
