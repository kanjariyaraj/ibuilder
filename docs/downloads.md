# Downloads

The download system retrieves artifacts from GitHub Actions workflow runs.

## Commands

### `ibuilder artifact download`

Download artifacts by various criteria.

```bash
# Download latest artifact
ibuilder artifact download --latest

# Download by name
ibuilder artifact download --name "iOS-Build"

# Download by build number
ibuilder artifact download --build 42

# Custom destination
ibuilder artifact download --latest --dest ./downloads

# Overwrite existing files
ibuilder artifact download --latest --overwrite
```

### `ibuilder artifact latest`

Quick download of the latest successful artifact.

```bash
ibuilder artifact latest
```

Saves to `dist/` directory automatically.

## Download Integrity

Downloads include:
- SHA256 checksum verification
- File size validation
- Error recovery on partial downloads

## Storage

Downloaded artifacts are stored locally with metadata:

```
.build/
├── artifacts/
├── metadata/
│   └── artifact-<id>.json
└── logs/
```
