# TestFlight Management

Builder provides complete TestFlight release management.

## Commands

- `builder testflight upload` — Upload IPA to TestFlight
- `builder testflight status` — Check upload/processing status
- `builder testflight groups` — List beta groups
- `builder testflight builds` — List builds
- `builder testflight testers` — List testers

## Upload Options

Upload latest IPA:
```
builder testflight upload
```

Upload specific artifact:
```
builder testflight upload --artifact path/to/app.ipa
```

Upload specific build:
```
builder testflight upload --build 42
```

## Status

Check overall status:
```
builder testflight status
```

Check specific build:
```
builder testflight status --build 42
```
