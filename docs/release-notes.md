# Release Notes Generation

Builder generates release notes from git history with AI-powered categorization.

## How It Works

1. Fetches the last 30 commits via `git log`
2. Categorizes commits into:
   - Features (feat, feature, add)
   - Bug Fixes (fix, bug, patch)
   - Breaking Changes (breaking, major)
3. Generates formatted output

## Output Formats

Markdown (default):
```
builder release notes
```

JSON:
```
builder release notes --format json
```

HTML:
```
builder release notes --format html
```
