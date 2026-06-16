# Configuration

Builder uses a `builder.json` file for configuration.

## Fields

| Field | Type | Description | Default |
|-------|------|-------------|---------|
| `project_name` | string | Name of the project | `"Builder"` |
| `repository` | string | Git repository URL | `""` |
| `ios.minimum_version` | string | Minimum iOS version | `"15.0"` |
| `ios.target_version` | string | Target iOS version | `"17.0"` |
| `ios.devices` | array | Target devices | `["iPhone", "iPad"]` |
| `signing.team_id` | string | Apple Team ID | `""` |
| `signing.provisioning_profile` | string | Provisioning profile path | `""` |
| `signing.certificate` | string | Signing certificate | `""` |
| `mob_ai.enabled` | bool | Enable MobAI | `false` |
| `mob_ai.api_key` | string | MobAI API key | `""` |
| `flutter.enabled` | bool | Enable Flutter | `false` |
| `flutter.channel` | string | Flutter channel | `"stable"` |
| `react_native.enabled` | bool | Enable React Native | `false` |
| `react_native.entry_file` | string | React Native entry file | `"index.js"` |

## Commands

```bash
# Create default configuration
builder config init

# Display current configuration
builder config show

# Validate configuration
builder config validate
```
