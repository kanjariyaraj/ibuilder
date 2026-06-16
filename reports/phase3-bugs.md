# Phase 3 Bug Report

## Fixed Issues

1. **APIRoot const cannot be reassigned** in tests
   - Changed from `const` to `var` to allow mock server URL injection

2. **Broken TestConstants** function syntax
   - Missing closing brace and incomplete if block fixed

3. **Missing NewRequest method** on Client
   - Added for POST/dispatch requests with JSON body

## Known Issues

- Workflow dispatch requires a valid workflow file in the repository
- Build logs are not streamed (URL provided for GitHub UI)
- Artifact download not fully implemented (listing only)
