# Flutter Development

## Overview

Build, deploy, and debug Flutter iOS apps on real devices via MobAI — no Mac required.

## Commands

### `builder flutter dev`

Start a development session: validate project, build, install, launch, and attach debugger.

```bash
builder flutter dev
builder flutter dev --device <device-id>
builder flutter dev --install
```

### `builder flutter attach`

Attach to an already-running Flutter app:

```bash
builder flutter attach
builder flutter attach --device <device-id>
```

### `builder flutter doctor`

Check Flutter development environment:

```bash
builder flutter doctor
```

Checks: Flutter SDK, Dart SDK, project validity, dependencies, device availability.

## Configuration

Flutter settings in `builder.json`:

```json
{
  "flutter": {
    "enabled": true,
    "channel": "stable",
    "watch": true,
    "hot_reload": true,
    "debounce_ms": 500,
    "auto_attach": true,
    "auto_install": true
  }
}
```
