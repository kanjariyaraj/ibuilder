# Artifacts

The artifact management system provides listing, downloading, inspecting, and cleaning of build artifacts from GitHub Actions.

## Commands

### `ibuilder artifact list`

List artifacts from workflow runs.

```bash
# List recent artifacts
ibuilder artifact list

# Limit results
ibuilder artifact list --limit 10

# Show all artifacts
ibuilder artifact list --all

# JSON output
ibuilder artifact list --json
```

### `ibuilder artifact download`

Download artifacts by various criteria.

```bash
# Download latest
ibuilder artifact download --latest

# Download by name
ibuilder artifact download --name "iOS-Build"

# Download by build number
ibuilder artifact download --build 42

# Custom destination
ibuilder artifact download --latest --dest ./output

# Overwrite existing
ibuilder artifact download --latest --overwrite
```

### `ibuilder artifact inspect`

View detailed artifact metadata.

```bash
ibuilder artifact inspect --id 12345
```

### `ibuilder artifact latest`

Download the latest successful artifact to `dist/`.

```bash
ibuilder artifact latest
```

### `ibuilder artifact clean`

Manage local artifact storage.

```bash
# Remove old artifacts (older than 7 days)
ibuilder artifact clean --older-than 168h

# Keep only recent artifacts
ibuilder artifact clean --keep 10

# Remove everything
ibuilder artifact clean --all
```

## Artifact Types

Builder detects and downloads:
- IPA files (iOS app packages)
- ZIP archives (build outputs)
- Any artifact uploaded by the workflow

## Build Report

The `build-report.json` contains run metadata:

```json
{
  "run_id": 1234567890,
  "run_number": 42,
  "status": "completed",
  "conclusion": "success",
  "duration": "120s",
  "workflow_url": "https://github.com/...",
  "artifact": "dist/build.zip"
}
```
