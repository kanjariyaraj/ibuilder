# iOS Build Engine

Build iOS apps remotely using GitHub Actions macOS runners.

## Commands

### Build

```bash
builder ios build [flags]
```

Start a remote iOS build via GitHub Actions.

### Flags

| Flag | Description |
|------|-------------|
| `--workflow` | Workflow ID or filename to dispatch |
| `--branch` | Branch to build from (default: from config) |
| `--scheme` | Xcode scheme name |
| `--mode` | Build mode: debug or release (default: release) |
| `--wait` | Wait for build to complete, then download artifacts |
| `--logs` | Stream build logs |
| `--json` | Output in JSON format |
| `--download-only` | Download latest artifact without triggering a new build |
| `--clean` | Clean the dist/ directory |

## Build Flow

1. Validate authentication
2. Validate repository configuration
3. Dispatch workflow via GitHub Actions API
4. (Optional) Wait for completion with live status updates
5. Download artifacts to `dist/`
6. Generate build report

## Examples

```bash
# Trigger a build
builder ios build --workflow ios-build.yml

# Build with custom branch and scheme
builder ios build --workflow ios-build.yml --branch develop --scheme MyApp

# Wait for completion and download artifacts
builder ios build --workflow ios-build.yml --wait

# Download latest successful artifact
builder ios build --download-only

# JSON output for CI
builder ios build --workflow ios-build.yml --json

# Clean dist/ directory
builder ios build --clean
```

## Output

Build artifacts are downloaded to the `dist/` directory:

```
dist/
├── <artifact-name>.zip    # Build artifact
├── build-report.json       # Machine-readable report
└── build-report.md         # Human-readable report
```
