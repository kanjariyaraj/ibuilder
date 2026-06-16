# Logs

Build logs provide detailed output from GitHub Actions workflow runs.

## Commands

### `ibuilder build logs`

Fetch and save build logs.

```bash
# Latest build logs
ibuilder build logs --latest

# Specific build
ibuilder build logs --run-id 1234567890

# Save to custom path
ibuilder build logs --latest --save ./my-build.log
```

## Storage

Logs are saved to `.build/logs/` by default with the naming convention `run-<id>.log`.

## Use Cases

- Debug build failures offline
- Share logs with team members
- Archive logs for compliance
- Parse logs with external tools
