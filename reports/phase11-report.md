# Phase 11 Report: AI Troubleshooter & Auto-Fix Engine

## Overview

Implemented AI-powered diagnostic engine for automatic build failure
analysis, root cause identification, and fix suggestions.

## Architecture Decisions

1. **Pattern-based classification**: Rule-based classifier system matches
   known failure patterns against log lines. More specific classifiers
   (signing, flutter, RN) take priority over generic build errors.

2. **Deterministic scoring**: Classifier ordering ensures reproducible
   results — first classifier with highest score wins ties.

3. **Knowledge base**: Built-in rule database covers all major failure
   categories with actionable fix suggestions.

4. **Safety-first fixes**: All fix operations are dry-run by default;
   explicit `--apply` flag required for changes.

5. **Report formats**: Support for markdown, JSON, and HTML output.

## Files Created

### internal/ai/ (8 files)
- `ai.go` — Session struct, log collection, project dir
- `analyzer.go` — Pattern classifiers, analysis engine
- `doctor.go` — System checks (Node, Go, Git, config, logs, deps)
- `explain.go` — Failure explanation with confidence
- `fix.go` — Fix suggestion generation, patch application
- `report.go` — Multi-format report generation
- `knowledgebase.go` — Rule database for all failure categories
- `classifiers.go` — Log classification and grouping

### cmd/builder/cmd/
- `ai.go` — CLI commands for `builder ai *`

### docs/ (4 files)
- `docs/ai-engine.md`, `docs/ai-doctor.md`, `docs/ai-fix.md`, `docs/troubleshooting.md`

### reports/
- `reports/phase11-report.md`
- `reports/phase11-testing.md`
- `reports/phase11-bugs.md`
- `reports/phase11-summary.md`
- `reports/phase11-completion.md`
