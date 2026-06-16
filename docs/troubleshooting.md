# Troubleshooting Guide

## AI Troubleshooter

### "No logs available"
Run a build first: `builder build` or `builder rn dev`

### False positives
Adjust confidence threshold in builder.json:
```json
"ai": {
  "confidence_threshold": 0.9
}
```

### Reports not generating
Check `.build/reports/ai/` directory exists and is writable.

## TestFlight

### Upload fails
Check:
1. IPA file exists in `.build/` directory
2. IPA file is valid (non-empty, correct format)
3. App Store Connect credentials are configured

### Status shows "PENDING"
Builds take time to process. Check back later with:
```
builder testflight status
```

### Groups not showing
Ensure TestFlight is enabled for your App Store Connect account.
