# Troubleshooting Authentication

## Common Issues

### "Not authenticated" Error

```bash
builder repo info
Error: not authenticated. Run 'builder auth github' first.
```

Run `builder auth github` to authenticate.

### Token Expired

If a token was previously valid but no longer works:

```bash
builder auth logout
builder auth github
```

### "Corrupt token file"

The token file at `~/.builder/github.json` may be corrupted.

**Recovery**:
```bash
builder auth logout
builder auth github
```

### Rate Limiting

GitHub API rate limits apply. If you receive rate limit errors:
- Wait before retrying
- Authenticated requests have higher limits

### No Internet Connection

Builder requires internet access to:
- Authenticate with GitHub
- Fetch repository information

Check your network connection and try again.

### Permission Denied

If you get permission errors:
- Verify your token has the required scopes
- Check repository access permissions
- Re-authenticate with `builder auth github`
