# Contributing to PortDoctor

Thank you for your interest in contributing!

## Good first issues

- Adding new framework detectors (Bun, Deno, Laravel, Rails, Spring Boot, etc.)
- Adding new runtime detectors
- Improving OS compatibility
- Writing tests
- Improving documentation
- Terminal rendering improvements

## Development

```bash
git clone https://github.com/UPin2905/portdoctor
cd portdoctor
go build ./cmd/portdoctor
go test ./...
```

## Before submitting

```bash
go fmt ./...
go vet ./...
go test ./...
go build ./cmd/portdoctor
```

All four commands must succeed.

## Adding a framework detector

Edit `internal/framework/detect.go`. Add pattern matching for the new framework and a corresponding test in `detect_test.go`.

## Adding an OS

Create `internal/port/port_<os>.go` and `internal/process/process_<os>.go` with the appropriate build tags.

## Code style

- Idiomatic Go
- Short functions
- Explicit error handling
- No unnecessary comments
- No global mutable state

## Pull requests

- Keep changes focused
- Include tests for new logic
- Update documentation if needed

