# Project Structure

```
Builder/
├── cmd/
│   └── builder/
│       ├── main.go          # Entry point
│       └── cmd/
│           ├── root.go      # Root command
│           ├── version.go   # Version command
│           ├── doctor.go    # Doctor command
│           └── config.go    # Config command
├── internal/
│   ├── config/              # Configuration management
│   ├── errors/              # Error handling
│   └── logger/              # Structured logging
├── docs/                    # Documentation
├── reports/                 # Reports and analysis
├── templates/               # Project templates
├── examples/                # Usage examples
├── tests/                   # Integration tests
├── roadmap/                 # Feature roadmap
├── .github/                 # GitHub workflows and templates
├── builder.json             # Default configuration
├── go.mod                   # Go module definition
└── go.sum                   # Go module checksum
```
