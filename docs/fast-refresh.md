# Fast Refresh & Reload

React Native supports two reload mechanisms.

## Fast Refresh

Triggered automatically on file save when Metro is running.
Can also be triggered manually:

```
builder rn reload --fast-refresh
```

Sends an HTTP POST to Metro's `/onchange` endpoint.

## Manual Reload

Full JavaScript bundle reload:

```
builder rn reload
```

Sends an HTTP POST to Metro's `/reload` endpoint.

## Requirements

- Metro bundler must be running
- React Native 0.61+ (Fast Refresh built in)
- `fast_refresh: true` in builder.json config
