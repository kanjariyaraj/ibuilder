# Storage

The local storage system manages build artifacts, logs, reports, and metadata.

## Directory Structure

```
.build/
├── artifacts/          # Downloaded artifact files
├── logs/               # Build logs
├── reports/            # Generated reports
├── cache/              # Temporary cache
└── metadata/           # Artifact and build metadata
```

## Metadata

Artifact metadata is stored as JSON files:

```json
{
  "id": 12345,
  "name": "iOS-Build",
  "local_path": "dist/iOS-Build.zip",
  "size": 1048576,
  "checksum": "abc123...",
  "downloaded_at": "2024-01-15T10:30:00Z"
}
```

## Cleanup

### `ibuilder artifact clean`

Remove old or all local artifacts.

```bash
# Remove artifacts older than 7 days
ibuilder artifact clean --older-than 168h

# Remove artifacts older than 30 days
ibuilder artifact clean --older-than 720h

# Keep only the 5 most recent artifacts
ibuilder artifact clean --keep 5

# Remove all local artifacts
ibuilder artifact clean --all
```
