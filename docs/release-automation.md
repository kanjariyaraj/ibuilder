# Release Automation

Builder provides end-to-end release automation from a single command.

## Workflow

```
ibuilder release --production
```

This single command:

1. **Validates** the environment and configuration
2. **Builds** the iOS app
3. **Verifies** code signing
4. **Generates** release notes
5. **Uploads** to TestFlight
6. **Creates** a GitHub release
7. **Reports** the results

## CI/CD Integration

The release pipeline can be integrated into CI/CD workflows:

```yaml
# .github/workflows/release.yml
steps:
  - uses: actions/checkout@v4
  - run: builder release --production --dry-run
```

## Safety

- All releases are **dry-run** by default for preview
- Sensitive data is masked in output
- Rollback infrastructure is prepared for future use
- Every step is logged for audit
