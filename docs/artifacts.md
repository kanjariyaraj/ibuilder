# Artifacts

Build artifacts are downloaded to the `dist/` directory.

## Artifact Types

Builder detects and downloads:

- IPA files (iOS app packages)
- ZIP archives (build outputs)
- Any artifact uploaded by the workflow

## Artifact Structure

```
dist/
├── <artifact-name>.zip
├── build-report.json
└── build-report.md
```

## Build Report

The `build-report.json` contains:

```json
{
  "run_id": 1234567890,
  "run_number": 42,
  "status": "completed",
  "conclusion": "success",
  "duration": "2024-01-15T10:30:00Z",
  "workflow_url": "https://github.com/...",
  "artifact": "dist/build.zip"
}
```

The `build-report.md` provides a human-readable summary.
