# Release Pipeline

Builder's one-command release pipeline automates the entire iOS release
process: validate → build → sign → notes → upload → GitHub release → report.

## Usage

Basic release:
```
builder release
```

Beta release:
```
builder release --beta
```

Production release:
```
builder release --production
```

Dry run (preview without changes):
```
builder release --dry-run
```

Release notes only:
```
builder release --notes
```

## Pipeline Stages

| # | Stage | Description |
|---|-------|-------------|
| 1 | Validate | Environment, config, signing, repo, GitHub, project |
| 2 | Build | Trigger workflow, build IPA, collect artifact |
| 3 | Sign | Verify certificate, provisioning, bundle signing |
| 4 | Notes | Generate release notes from git history |
| 5 | Upload | Upload IPA to TestFlight |
| 6 | Release | Create GitHub release with IPA and notes |
| 7 | Report | Generate release report and audit log |

## Release Modes

- **beta** — TestFlight beta distribution
- **production** — Production App Store release
- **internal** — Internal testing release
- **custom** — Custom release configuration

## Status

Check release pipeline status:
```
builder release status
```
