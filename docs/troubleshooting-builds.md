# Troubleshooting Builds

## Common Issues

### "not authenticated"

```bash
builder auth github
```

### "no repository configured"

```bash
builder repo connect
```

### "workflow ID required"

Set the workflow ID in `builder.json`:
```json
{
  "build": {
    "workflow_id": "ios-build.yml"
  }
}
```

Or pass it via CLI:
```bash
builder ios build --workflow ios-build.yml
```

### Build fails on GitHub

Check the workflow run URL for detailed logs:
```bash
builder build status <run-id>
```

### Artifact not found

- Ensure your workflow uploads artifacts using `actions/upload-artifact`
- Check the artifact name in your workflow file
- Verify the build completed successfully
