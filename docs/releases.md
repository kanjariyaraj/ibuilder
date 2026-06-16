# Release Management

Builder handles the full release lifecycle from notes generation
to release validation.

## Commands

- `builder release notes` — Generate release notes from git history
- `builder release history` — View release history
- `builder release prepare` — Validate release readiness

## Release Notes

Generate with automatic categorization of commits:

```
builder release notes
builder release notes --format json
builder release notes --format html
```

Notes are saved to `.build/releases/`.

## Release Preparation

Before releasing, validate everything:

```
builder release prepare
```

Checks:
- Signing configuration
- Build artifacts
- IPA file validity
- App metadata
- Release notes
- Git state

## Release History

View all releases:
```
builder release history
```

Inspect a specific version:
```
builder release history --version 1.0.0
```
