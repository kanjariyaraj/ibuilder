# AI Troubleshooter Engine

Builder's AI engine provides intelligent diagnostics for build failures,
workflow issues, signing problems, and more.

## Commands

- `builder ai explain` — Explain the latest failure
- `builder ai analyze` — Deep-dive analysis of logs
- `builder ai doctor` — Full repository audit
- `builder ai fix` — Generate fix suggestions
- `builder ai report` — Generate diagnostic reports

## How It Works

The AI engine uses pattern-based classification to categorize failures:

1. **Collect logs** from `.build/logs/` directory
2. **Classify** each log line against known failure patterns
3. **Score** categories by match frequency
4. **Determine** root cause with confidence level
5. **Suggest** fixes from the knowledge base

## Failure Categories

- BUILD — Xcode compilation errors
- SIGNING — Code signing and provisioning issues
- FLUTTER — Flutter/Dart build failures
- REACT_NATIVE — React Native and Metro errors
- DEPENDENCY — Missing or conflicting dependencies
- NETWORK — Connectivity issues
- PERMISSION — Access and authorization problems
- GITHUB_ACTIONS — CI/CD workflow failures

## Confidence Levels

- **HIGH** (≥80%): Strong match with multiple corroborating signals
- **MEDIUM** (50-79%): Plausible match with some supporting evidence
- **LOW** (<50%): Weak match, may need manual investigation

## Safety

The AI engine never applies destructive changes without explicit
user approval. All fixes are dry-run by default.
