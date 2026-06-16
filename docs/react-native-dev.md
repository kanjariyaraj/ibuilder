# React Native Development

Builder provides a complete React Native development workflow for iOS apps
on real devices through MobAI, without needing a Mac.

## Commands

- `builder rn dev` — Start development mode
- `builder rn attach` — Attach to a running app
- `builder rn metro` — Manage Metro bundler
- `builder rn reload` — Trigger Fast Refresh or manual reload
- `builder rn logs` — View and stream device logs
- `builder rn doctor` — Check development environment
- `builder rn install` — Install app on device

## Prerequisites

- Node.js 18+
- npm
- React Native project with `ios/` directory
- MobAI device connection (optional but recommended)

## Configuration

Add to `builder.json`:

```json
{
  "react_native": {
    "enabled": true,
    "entry_file": "index.js",
    "metro_port": 8081,
    "auto_start_metro": true,
    "auto_attach": true,
    "auto_install": true,
    "fast_refresh": true
  }
}
```

## Workflow

1. `builder rn doctor` — verify environment
2. `builder rn dev` — start development session
3. Edit your code — Metro hot reloads automatically
4. `builder rn reload` — manual reload if needed
5. `builder rn logs` — view device logs
6. `builder rn attach` — attach to existing session
