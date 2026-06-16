# GitHub Releases

Builder integrates with GitHub Releases to attach build artifacts,
release notes, and reports.

## Automatic GitHub Release

When running the release pipeline, Builder automatically:

1. Creates a GitHub release with version tag
2. Attaches the IPA artifact
3. Includes generated release notes
4. Attaches diagnostic reports

## Configuration

Enable/disable in builder.json:
```json
{
  "release": {
    "create_github_release": true
  }
}
```

## Tagging

Release tags follow the format: `vYYYY.MM.DD`
