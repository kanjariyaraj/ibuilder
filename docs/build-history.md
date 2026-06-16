# Build History

The build history system provides access to past workflow runs and their details.

## Commands

### `ibuilder build history`

List recent builds with filtering options.

```bash
# Show last 30 builds
ibuilder build history

# Filter by branch
ibuilder build history --branch main

# Filter by status
ibuilder build history --status completed

# Limit results
ibuilder build history --limit 10

# Pagination
ibuilder build history --page 2 --limit 20

# JSON output
ibuilder build history --json
```

### `ibuilder build inspect`

Show detailed information about a specific build.

```bash
# Inspect by run ID
ibuilder build inspect --run-id 1234567890
```

Displays:
- Build number and status
- Branch and commit SHA
- Author and workflow name
- Duration and URL
- Artifacts produced
- Job breakdown with step counts

### `ibuilder build logs`

Download logs for a build.

```bash
# Latest build logs
ibuilder build logs --latest

# Specific build logs
ibuilder build logs --run-id 1234567890

# Custom save path
ibuilder build logs --latest --save ./logs/build.log
```

Logs are saved to `.build/logs/` by default.

### `ibuilder build open`

Open a build's GitHub Actions page in your browser.

```bash
ibuilder build open --run-id 1234567890
```

## Output Format

Build history displays in a table:

```
#      Status      Conclusion   Branch               Duration
42     completed   success      main                 120s
41     completed   failure      feature/new-ui       85s
```
