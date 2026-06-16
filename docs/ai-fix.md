# AI Fix

The AI fix engine analyzes failures and generates safe fix suggestions.

## Usage

Preview fixes (default):
```
builder ai fix
```

Apply fixes:
```
builder ai fix --apply
```

## Safety

- Fixes are never applied without confirmation
- Use `--dry-run` (default) to preview changes
- The `--apply` flag applies suggested fixes
- No destructive operations are performed

## What Can Be Fixed

- Build mode configuration
- Signing team ID setup
- Dependency installation
- Flutter configuration
- React Native configuration
- Network connectivity checks
