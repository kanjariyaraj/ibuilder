# Phase 10 Report: React Native Development Workflow

## Overview

Implemented complete React Native development workflow for Windows, Linux,
and WSL users using real iPhone devices through MobAI.

## Architecture Decisions

1. **Package naming**: `internal/reactnative/` follows snake_case Go convention,
   consistent with existing `internal/flutter/`.

2. **Session model**: Mirror Flutter's Session pattern with mutex-guarded state,
   supporting Inactive → Starting → Active → Attached → MetroRunning states.

3. **Metro management**: Separate metro.go module handles start/stop/restart/status
   lifecycle, with port conflict detection via lsof.

4. **Reload mechanism**: HTTP-based Fast Refresh and manual reload through Metro's
   `/onchange` and `/reload` endpoints.

5. **Config structure**: Extended `ReactNativeSettings` in config.go with metro_port,
   auto_start_metro, auto_attach, auto_install, fast_refresh fields.

## Files Created

### internal/reactnative/ (10 files)
- `reactnative.go` — Session struct, NewSession, project detection, Node checks
- `doctor.go` — Health checking (Node, npm, project, deps, Metro, devices)
- `dev.go` — Dev mode with auto Metro start
- `attach.go` — Attach to running RN app
- `metro.go` — Metro bundler lifecycle management
- `reload.go` — Fast Refresh and manual reload
- `logs.go` — Log fetching, streaming, filtering, saving
- `install.go` — App installation (latest and custom artifact)
- `session.go` — Session info, stop, active check
- `recovery.go` — Auto recovery up to 3 attempts

### cmd/builder/cmd/
- `reactnative.go` — All CLI commands for `builder rn *`

### docs/
- `docs/react-native-dev.md`
- `docs/metro.md`
- `docs/fast-refresh.md`
- `docs/react-native-troubleshooting.md`

### reports/
- `reports/phase10-report.md`
- `reports/phase10-testing.md`
- `reports/phase10-bugs.md`
- `reports/phase10-summary.md`
- `reports/phase10-completion.md`
