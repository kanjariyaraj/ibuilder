# Beta Testing

Manage TestFlight beta testing groups and testers.

## Commands

- `builder testflight groups` — List and inspect groups
- `builder testflight testers` — List testers

## Groups

List all groups:
```
builder testflight groups
```

Inspect a specific group:
```
builder testflight groups --inspect "Internal Testers"
```

## Testers

List all testers:
```
builder testflight testers
```

## Build Distribution

Upload builds to TestFlight:
```
builder testflight upload
```

Check build availability:
```
builder testflight status
```
