# Phase 11 Bugs & Fixes

## Compilation Issues

1. **Missing `os` import in classifiers.go**
   - Error: `undefined: os`
   - Fix: Added `"os"` to imports

2. **Analyzer.knowledgeBase undefined**
   - Error: `a.knowledgeBase undefined (type *Analyzer has no field or method knowledgeBase)`
   - Fix: Changed to create local KnowledgeBase instance in Analyze method

## Logic Issues

3. **Non-deterministic classifier results**
   - Problem: Go map iteration is randomized, causing test failures
   - Fix: Added deterministic `classifierOrder` slice for tie-breaking
   - Specific classifiers (network, signing) now take priority over build

4. **Nil logger panic in tests**
   - Problem: Tests pass nil logger, causing SIGSEGV on log.Info calls
   - Fix: Added nil-safe `logInfo`/`logWarn` wrapper methods

## No Other Test Failures

All tests pass on first run after fixes.
