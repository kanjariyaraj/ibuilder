# iOS Build Pipeline

Builder can trigger iOS builds on GitHub Actions macOS runners.

## Prerequisites

1. GitHub authentication configured (`builder auth github`)
2. Repository connected (`builder repo connect`)
3. A GitHub Actions workflow for iOS builds in your repository

## Commands

### Trigger a Build

```bash
builder build run <workflow-id> [flags]
```

Flags:
- `--branch` - Branch to build from (default: main)
- `--scheme` - Xcode scheme name
- `--mode` - Build mode (debug/release, default: release)

Example:
```bash
builder build run ios-build.yml --branch main --scheme MyApp --mode release
```

### Check Build Status

```bash
builder build status <run-id>
```

### List Recent Builds

```bash
builder build list --limit 20
```

### View Build Logs

```bash
builder build log <run-id>
```

### List Build Artifacts

```bash
builder build artifacts <run-id>
```

## Workflow Templates

Builder includes workflow templates in the `templates/` directory:

| Template | Description |
|----------|-------------|
| `ios-xcode.yml` | Standard Xcode iOS build |
| `ios-flutter.yml` | Flutter iOS build |
| `ios-react-native.yml` | React Native iOS build |

Copy a template to your repository's `.github/workflows/` directory to get started.
