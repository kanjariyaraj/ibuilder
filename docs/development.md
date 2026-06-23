# Development

## Prerequisites

- Go 1.21+

## Setup

```bash
git clone https://github.com/kanjariyaraj/iBuilder.git
cd Builder
go mod tidy
```

## Using the Makefile

This project includes a `Makefile` for streamlined development tasks.

```bash
# Build the project
make build

# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run linters
make lint

# Run all checks (lint, test, build)
make all
```

## Manual Commands

If you prefer not to use `make`, you can run the Go commands directly:

### Build

```bash
go build -o builder ./cmd/builder/
```

### Test

```bash
go test ./... -v
```

### Lint

```bash
go vet ./...
```

## Project Conventions

- **SOLID principles** — Single responsibility, open/closed, Liskov substitution, interface segregation, dependency inversion
- **Clean Architecture** — Separation of concerns, dependency rule
- **Small packages** — Each package has a single responsibility
- **Reusable code** — No duplication, extract common patterns
- **No dead code** — Remove unused functions and variables
- **No unused imports** — Run `go vet` before committing
- **No debug leftovers** — Remove debug statements
- **Production ready only** — Only production-quality code

## Commit Convention

```
feat(scope): description

Examples:
feat(core): initialize Builder foundation architecture and CLI framework
feat(config): add configuration management
feat(doctor): add system dependency checker
```
