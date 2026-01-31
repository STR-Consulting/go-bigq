---
description: Critical code review of uncommitted changes. Runs lints, tests, and reviews code for issues.
---

# Critical Code Review

## Step 1: Check for Changes

Run `git status --short` to verify there are changes to review.

If no changes, report "Nothing to review" and stop.

## Step 2: Run Lints

Run `golangci-lint run ./...`

If lint errors remain:
- Report each error with file:line
- These MUST be fixed before proceeding

## Step 3: Run Go Tests

Run `go test -v ./...`

If tests fail:
- Report failing tests with error messages
- These MUST be fixed before proceeding

## Step 4: Critical Code Review

Read the diff: `git diff HEAD`

Evaluate against this checklist:

### Code Quality
- [ ] Error handling is complete (no ignored errors)
- [ ] No magic strings/numbers (use constants)
- [ ] No unused code/imports
- [ ] Functions are reasonably sized
- [ ] No obvious performance issues

### Security
- [ ] No hardcoded secrets/credentials
- [ ] No command injection risks

### Tests
- [ ] New code has test coverage
- [ ] Tests are meaningful (not just coverage padding)

## Step 5: Report Findings

Summarize:
- Passed checks
- Warnings (non-blocking)
- Blocking issues that must be fixed
