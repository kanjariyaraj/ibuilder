# Repository Management

Builder can connect to and manage GitHub repositories.

## Commands

### Connect Repository

Detect the repository from the current git remote:

```bash
builder repo connect
```

This reads the `origin` remote URL and saves the owner and name to `builder.json`.

### Repository Info

Display metadata about the connected repository:

```bash
builder repo info
```

Output includes:
- Owner and repository name
- Default branch
- Visibility (public/private)
- User permissions (admin, push, pull)
- GitHub Actions status

### Validate Repository

Verify repository access and configuration:

```bash
builder repo validate
```

Checks:
- Authentication is valid
- Repository exists and is accessible
- GitHub Actions are enabled
- User has sufficient permissions
