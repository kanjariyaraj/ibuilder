# Phase 12 Report: TestFlight & Release Management

## Overview

Implemented complete TestFlight upload and release management system.

## Architecture Decisions

1. **Session model**: Consistent with Flutter/RN/AI patterns — mutex-guarded
   Session with project directory management.

2. **IPA discovery**: Recursive walk of `.build/` directory to find IPA files
   without hardcoded paths.

3. **Release notes**: Git-based automatic categorization of commits into
   features, fixes, and breaking changes.

4. **Preparation checks**: Pre-flight validation of signing, builds, IPA,
   metadata, notes, and git state before release.

5. **Multiple formats**: Release notes and reports support markdown, JSON,
   and HTML output.

## Files Created

### internal/release/ (8 files)
- `release.go` — Session struct, IPA discovery
- `upload.go` — TestFlight upload (latest, artifact, build)
- `status.go` — Status checking (upload, processing, beta, review)
- `groups.go` — Beta group listing and inspection
- `builds.go` — Build listing and details
- `testers.go` — Tester listing
- `notes.go` — Release notes generation and formatting
- `history.go` — Release history tracking
- `prepare.go` — Pre-release validation

### cmd/builder/cmd/
- `release.go` — CLI commands for `builder testflight *` and `builder release *`

### docs/ (4 files)
- `docs/testflight.md`, `docs/releases.md`, `docs/release-notes.md`, `docs/beta-testing.md`

### reports/
- `reports/phase12-report.md`
- `reports/phase12-testing.md`
- `reports/phase12-bugs.md`
- `reports/phase12-summary.md`
- `reports/phase12-completion.md`
