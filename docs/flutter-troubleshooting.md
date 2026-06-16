# Flutter Troubleshooting

## Common Issues

### Flutter SDK Not Found

```
Error: Flutter SDK not found
```

Run `flutter doctor` or install Flutter from https://flutter.dev.

### No pubspec.yaml

```
Error: no pubspec.yaml found
```

Run `builder flutter dev` from a valid Flutter project directory.

### No iOS Directory

```
Error: no ios/ directory found
```

Run `flutter create --platforms=ios .` in your project.

### No Devices Found

```
Error: no devices found
```

Connect a device via `builder mobai connect` or start an iOS simulator.

### Hot Reload Fails

```
Error: no active flutter session
```

Start a session first: `builder flutter dev` or `builder flutter attach`.

### Dependency Resolution Fails

Run `flutter pub get` manually then retry.
