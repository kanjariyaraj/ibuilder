# Build Engine

The Builder Build Engine (`internal/build`) is the core orchestrator for remote iOS compilation.

## How It Works

1. **Workflow Dispatch**: The engine uses the GitHub REST API to trigger specialized workflows located in `.github/workflows/`.
2. **Dynamic Inputs**: It passes project-specific parameters (scheme, configuration, build mode) as workflow inputs.
3. **Real-time Monitoring**: Once triggered, the engine polls the GitHub API to monitor the build's progress.
4. **Log Streaming**: It retrieves and displays logs from the remote runner in real-time, providing a "local-feel" experience.
5. **Artifact Retrieval**: Upon successful completion, the engine automatically downloads the resulting `.ipa` files and build reports.

## Configuration

Build settings are managed in `builder.json` under the `build` and `ios` keys.

```json
{
  "build": {
    "workflow_id": "ios-build.yml",
    "branch": "main",
    "scheme": "MyApp",
    "configuration": "Release"
  }
}
```

## Commands

- `builder build run`: Start a new remote build.
- `builder build status`: Check the status of the current or most recent build.
- `builder build log`: View logs from the remote runner.
- `builder build artifacts`: List and download build outputs.
