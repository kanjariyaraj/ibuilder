# GitHub Authentication

Builder uses the GitHub Device Authorization Flow to authenticate with GitHub.

## Authentication Flow

1. Run `builder auth github`
2. A device code is displayed along with a verification URL
3. Open the URL in a browser (opens automatically if possible)
4. Enter the device code
5. Authorize the application
6. Builder automatically stores the token

## Commands

### Authenticate

```bash
builder auth github
```

### Check Status

```bash
builder auth status
```

### Logout

```bash
builder auth logout
```

## Token Storage

Tokens are stored securely at:

- **Linux/macOS**: `~/.builder/github.json`
- **Windows**: `%USERPROFILE%\.builder\github.json`

File permissions are set to `0600` (owner read/write only).

## Scopes Requested

- `repo` - Access to public and private repositories
- `workflow` - Access to GitHub Actions workflows
