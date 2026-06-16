# MobAI Protocol

MobAI is a specialized, high-performance interaction layer designed for remote iOS device management and debugging.

## Overview

MobAI solves the problem of interacting with physical iOS devices or simulators that are not locally connected. It establishes a secure TCP tunnel between the Builder CLI and a remote "Agent" running on a macOS machine.

## Key Features

- **Remote Hot Reload**: Seamlessly inject code changes into a running Flutter or React Native app on a remote device.
- **Device Diagnostics**: Retrieve real-time battery levels, storage status, and OS information from remote hardware.
- **Log Forwarding**: Stream system and application logs from the remote device directly to your local terminal.
- **Connection Persistence**: Automatic reconnection logic handles network fluctuations without interrupting your session.

## Architecture

The MobAI client (`internal/mobai`) manages the lifecycle of the connection:

1. **Discovery**: Identifies available remote agents.
2. **Handshake**: Establishes a secure session and verifies protocol versions.
3. **Heartbeat**: Maintains the connection and monitors latency.
4. **Data Transfer**: Handles the bidirectional flow of debugging commands and device data.

## Usage

Enable MobAI in your `builder.json`:

```json
{
  "mobai": {
    "host": "remote-mac.local",
    "port": 12345,
    "device": "iPhone 15 Pro",
    "auto_reconnect": true
  }
}
```
