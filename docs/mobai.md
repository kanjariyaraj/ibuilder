# MobAI Device Integration

## Overview

MobAI enables real iPhone device management from the CLI. Connect, manage, and debug iOS devices remotely.

## Commands

### `builder mobai connect`

Connect to a MobAI device:

```bash
builder mobai connect
builder mobai connect --host 192.168.1.100 --port 12345 --device "iPhone 15"
```

### `builder mobai disconnect`

Disconnect from the current device.

### `builder mobai status`

Display connection status, device info, and latency.

### `builder mobai doctor`

Run health checks: configuration, connectivity, device availability.

### `builder mobai ping`

Ping the device and measure latency.

## Configuration

MobAI settings in `builder.json`:

```json
{
  "mobai": {
    "host": "",
    "port": 0,
    "device": "",
    "auto_reconnect": true,
    "connection_timeout": 30
  }
}
```
