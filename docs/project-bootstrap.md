# Project Bootstrap

Quickly bootstrap any iOS project for use with Builder.

## Quick Start

```bash
# 1. Initialize project
cd your-project
builder init run

# 2. Authenticate with GitHub
builder auth github

# 3. Connect repository
builder repo connect

# 4. Trigger a build
builder build run ios-build.yml
```

## CI Integration

For CI environments:

```bash
builder init run --yes --json
```

## What Gets Created

After `builder init run`:

```
your-project/
├── builder.json           # Project configuration
└── .github/
    └── workflows/
        └── ios-build.yml  # GitHub Actions workflow
```
