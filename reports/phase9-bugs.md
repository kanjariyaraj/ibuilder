# Phase 9 Bug Report

## Fixed During Development

None — all tests passed on first run.

## Known Issues

1. Flutter SDK must be installed and in PATH for dev/attach to work
2. File watcher uses polling (500ms) — not as efficient as inotify/FSEvents
3. Mock project detection — real Flutter projects need `ios/` directory
4. Recovery requires Flutter SDK and a valid project
5. Log streaming uses process piping — may have buffering delays
