# Getting Started

## Prerequisites

- Go 1.21+
- Git
- GitHub CLI (optional, for GitHub integration)

## Installation

```bash
go install github.com/kanjariyaraj/Builder/cmd/builder@latest
```

## Quick Start

1. Create a default configuration:
   ```bash
   builder config init
   ```

2. Verify your system is ready:
   ```bash
   builder doctor
   ```

3. View your configuration:
   ```bash
   builder config show
   ```

## Building from Source

```bash
git clone https://github.com/kanjariyaraj/Builder.git
cd Builder
go build -o builder ./cmd/builder/
./builder version
```
