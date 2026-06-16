# Workflow Generation

Builder generates GitHub Actions workflow files for iOS builds.

## Generated Workflow

When you run `builder init run`, a workflow is created at `.github/workflows/ios-build.yml`.

Each project type gets a tailored workflow template:

| Project Type | Workflow Features |
|-------------|-------------------|
| Native iOS | Xcode build with scheme selection, Debug/Release config |
| Flutter | Flutter build with channel support, debug/release modes |
| React Native | Metro bundle with platform targeting |
| Expo | Expo build with profile configuration |
| Capacitor | Capacitor sync and build |
| Unity | Unity iOS build via game-ci |
| Unreal | Manual iOS build setup guide |

## Template System

Templates are stored in `templates/` and selected automatically based on project type:

- `ios-xcode.yml` → Native iOS projects
- `ios-flutter.yml` → Flutter projects
- `ios-react-native.yml` → React Native projects

## Workflow Inputs

Generated workflows accept inputs via `workflow_dispatch`:

```bash
builder build run ios-build.yml --branch main --scheme MyApp --mode release
```
