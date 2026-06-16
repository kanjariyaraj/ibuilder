# Device Logs

## Commands

### `builder device logs`

View device logs:

```bash
builder device logs                     # Show all logs
builder device logs --level ERROR       # Filter by level
builder device logs --process app       # Filter by process
builder device logs --search "timeout"  # Search logs
builder device logs --since 5m          # Last 5 minutes
builder device logs --save logs/        # Save to directory
builder device logs --stream            # Live streaming
```

## Output Location

Saved logs: `.build/logs/device_logs_<timestamp>.txt`
