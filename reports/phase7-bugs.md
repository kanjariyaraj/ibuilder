# Phase 7 Bug Report

## Fixed Issues

1. **parseDuration undefined in artifacts test** - The `internal/artifacts/artifacts_test.go` referenced `parseDuration` which was defined in `cmd/builder/cmd/artifact.go`, not the artifacts package. Fixed by using `time.ParseDuration` directly.

2. **Build tests referencing removed commands** - `cmd/builder/cmd/build_test.go` tested old commands (`run`, `status`, `list`, `log`, `artifacts`) that were replaced with new ones (`history`, `inspect`, `logs`, `open`). Updated tests to match the new command structure.

## Known Issues

1. **GitHub API dependency** - Download, history, inspect, and logs commands require valid GitHub authentication and repository configuration. Errors surface as "not authenticated" or API failure messages.

2. **Download URL expiration** - GitHub artifact download URLs expire after a limited time. The system uses the GitHub API to get fresh URLs.

3. **No automatic retry** - Network failures during download do not automatically retry. Users must re-run the command.

4. **No progress bar** - Large downloads show no progress indicator beyond completion message.
