# Architecture Plan

## Directory Responsibilities

### `cmd/builder/`
- Single main.go entry point
- Cobra root command setup
- Subcommand registration

### `internal/config/`
- Configuration schema and types
- Load/Save/Validate operations
- Default values
- JSON serialization

### `internal/logger/`
- Structured logging (Info, Warn, Error, Debug)
- Cross-platform output
- Consistent formatting

### `internal/errors/`
- Reusable error types
- Error wrapping with context
- User-friendly messages

### `pkg/` (future)
- Public API packages for extensibility
- Shared utilities

### `docs/`
- User-facing documentation
- Getting started, configuration, development guides

### `reports/`
- Architecture decisions
- Testing reports
- Bug reports
- Implementation summaries

### `templates/`
- Project templates for supported frameworks (Flutter, React Native, etc.)

### `examples/`
- Usage examples and sample configurations

### `tests/`
- Integration and end-to-end tests

### `roadmap/`
- Feature roadmap and version planning

### `.github/`
- GitHub Actions workflows
- Issue templates
- PR templates

## Architecture Decisions

1. **CLI Framework**: `github.com/spf13/cobra` - industry standard for Go CLIs
2. **Config Format**: JSON - human-readable, widely supported
3. **Logging**: Custom structured logger with levels
4. **Error Handling**: Wrapped errors with context using `fmt.Errorf("%w")`
5. **Testing**: Standard `testing` package with table-driven tests

## Component Relationships

```
cmd/builder/main.go
    └── root.go (cobra.Command)
        ├── version.go → outputs build info
        ├── doctor.go  → runs system checks via internal/doctor
        └── config.go  → manages config via internal/config
```

```
internal/config/
    ├── config.go     → Config struct, Load, Save, Validate
    └── defaults.go   → Default configuration values

internal/logger/
    └── logger.go     → Logger interface and implementation

internal/errors/
    └── errors.go     → Reusable error types and helpers
```
