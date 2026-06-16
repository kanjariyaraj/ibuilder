# Flutter Hot Reload & Restart

## Overview

Update your Flutter app instantly without rebuilding.

## Commands

### Hot Reload

```bash
builder flutter reload
```

Injects updated Dart code into the running app. Preserves app state.

### Hot Restart

```bash
builder flutter restart
```

Restarts the Dart VM and re-runs `main()`. Resets app state.

## Requirements

- Active Flutter session (`builder flutter dev` or `builder flutter attach`)
- Hot reload must be enabled in config

## Configuration

```json
{
  "flutter": {
    "hot_reload": true
  }
}
```
