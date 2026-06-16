# Release Modes

Builder supports multiple release modes for different distribution channels.

## Beta Mode

Distribute to TestFlight beta testers:

```
builder release --beta
```

- Uploads to TestFlight
- Available to internal testers
- No App Store review required for beta

## Production Mode

Prepare for App Store submission:

```
builder release --production
```

- Full validation of signing and metadata
- Upload for App Store review
- GitHub release with production tag

## Dry Run Mode

Preview all release steps without making changes:

```
builder release --dry-run
```

Useful for:
- CI pipeline validation
- Testing configuration
- Pre-release checks

## Notes Only

Generate release notes without running the full pipeline:

```
builder release --notes
```
