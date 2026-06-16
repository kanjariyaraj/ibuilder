# AI Doctor

The AI doctor performs a comprehensive audit of your project.

## What It Checks

| Check | Description |
|-------|-------------|
| Node.js | Is Node.js installed? |
| Go | Is Go available? |
| Git | Is Git installed? |
| Project Config | Is builder.json present? |
| Build Logs | Are build logs available with issues? |
| Dependencies | Are dependencies installed? |

## Usage

```
builder ai doctor
```

## Output

Each check shows a status:
- ✓ HEALTHY — No issues
- ⚠ WARNING — Minor concern (non-blocking)
- ✗ FAILURE — Blocking issue

The overall status is shown at the end.
