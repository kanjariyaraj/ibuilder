# Phase 8 Bug Report

## Fixed During Development

1. **Sanitize filename test**: Expected `Raj_s_iPhone` but apostrophes are skipped (not replaced). Fixed expectation to `Rajs_iPhone`.

## Known Issues

1. Mock data is used for devices/logs — real MobAI server integration will replace mocks
2. Screenshot generates a synthetic PNG — real device screenshots require MobAI API
3. Connection uses TCP dial — real deployment uses MobAI's HTTP/WebSocket API
4. Latency measurements are simulated
