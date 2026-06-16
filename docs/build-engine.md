# Build Engine Architecture

## Components

### `internal/build/`

| File | Responsibility |
|------|---------------|
| `runner.go` | Main build flow: validate, dispatch, orchestrate |
| `tracker.go` | Wait for build completion with status polling |
| `download.go` | Download artifacts, download-only mode |
| `report.go` | Generate build reports (JSON + Markdown) |

### Build Flow

```
builder ios build
  → Validate auth, repo, workflow
  → Dispatch workflow_dispatch event
  → (optional) Poll status every 10s
  → (on success) Download artifact to dist/
  → Generate build-report.json and build-report.md
```

## Configuration

Build settings read from `builder.json`:

```json
{
  "build": {
    "workflow_id": "ios-build.yml",
    "branch": "main",
    "scheme": "MyApp",
    "build_mode": "release"
  }
}
```

CLI flags override config values.
