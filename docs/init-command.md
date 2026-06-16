# Init Command

The `builder init run` command automatically detects your project type and generates the necessary configuration and GitHub Actions workflow.

## Usage

```bash
builder init run [directory]
```

If no directory is specified, the current directory is used.

## Flags

| Flag | Description |
|------|-------------|
| `--force` | Overwrite existing files without prompting |
| `--dry-run` | Show what would be done without making changes |
| `--yes` | Answer yes to all prompts (for CI automation) |
| `--json` | Output results in JSON format |

## Detection Flow

1. Detect project type (Flutter, React Native, Expo, Capacitor, Native iOS, Unity, Unreal)
2. Detect project name from configuration files
3. Detect git repository information
4. Detect iOS path, workspace, and scheme
5. Generate `builder.json` with auto-populated values
6. Generate `.github/workflows/ios-build.yml` with the appropriate template

## Supported Project Types

| Type | Detected By |
|------|-------------|
| Flutter | `pubspec.yaml` |
| React Native | `package.json` with `react-native` dependency |
| Expo | `package.json` with `expo` dependency |
| Capacitor | `package.json` with `@capacitor` dependency |
| Native iOS | `.xcodeproj` or `.xcworkspace` |
| Unity | `Assets/` and `ProjectSettings/` directories |
| Unreal | `.uproject` file or `Source/` and `Config/` directories |

## Examples

```bash
# Initialize current directory
builder init run

# Dry run to see what will be created
builder init run --dry-run

# Force overwrite existing files
builder init run --force

# CI mode - no prompts
builder init run --yes

# JSON output
builder init run --json

# Specify a project directory
builder init run /path/to/my-project
```
