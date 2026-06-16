# Troubleshooting

This guide provides solutions to common issues encountered when using Builder.

## 1. Authentication Issues

### "Not authenticated with GitHub"
- **Cause**: Your session has expired or you haven't logged in yet.
- **Fix**: Run `builder auth github` and follow the device flow instructions.

### "Permission denied (publickey)"
- **Cause**: The remote repository requires SSH keys that are not configured on the remote runner.
- **Fix**: Ensure your GitHub Token has the `repo` scope or use HTTPS URLs for your repository.

## 2. Build Failures

### Code Signing Errors
- **Cause**: Missing or invalid Provisioning Profiles or Certificates.
- **Fix**: Check your `builder.json` signing section. Use `builder doctor` to verify your local configuration.

### "Scheme not found"
- **Cause**: The specified Xcode scheme is not shared or does not exist.
- **Fix**: In Xcode, go to **Product > Scheme > Manage Schemes** and ensure the "Shared" checkbox is checked for your target scheme.

## 3. MobAI Connectivity

### "Connection timed out"
- **Cause**: The remote agent is not reachable or the port is blocked.
- **Fix**: Verify the `host` and `port` in `builder.json`. Ensure your firewall allows TCP traffic on the specified port.

## 🧠 Using the AI Doctor

If you're stuck, use the built-in AI diagnostics:

```bash
builder ai fix
```

The AI Doctor will scan your logs and configuration to provide a precise solution tailored to your environment.
