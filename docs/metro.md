# Metro Bundler Management

Builder provides full Metro bundler lifecycle management.

## Commands

- `builder rn metro start` — Start Metro on default port 8081
- `builder rn metro stop` — Stop Metro
- `builder rn metro restart` — Restart Metro
- `builder rn metro status` — Check if Metro is running

## Options

Start on a custom port:
```
builder rn metro start --port 8082
```

## Port Conflicts

Builder detects port conflicts automatically during dev mode.
If port 8081 is in use, check with:

```
builder rn metro status
lsof -i :8081
```

## Recovery

Metro crashes and disconnections are handled automatically
with the recovery system (up to 3 retry attempts).
