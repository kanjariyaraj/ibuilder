# Flutter Watch Mode

## Overview

Watch mode monitors file changes and automatically triggers hot reloads.

## Command

```bash
builder flutter watch
```

## Features

- File system monitoring for `.dart`, `.yaml`, `.xml`, `.plist`, `.json` files
- Smart debounce (configurable, default 500ms)
- Ignores generated directories (`.dart_tool`, `build`, `.git`, `Pods`, etc.)
- Detects created, modified, and deleted files
- Displays change events in real-time

## Configuration

```json
{
  "flutter": {
    "watch": true,
    "debounce_ms": 500
  }
}
```
