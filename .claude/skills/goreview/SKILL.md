---
name: goreview
description: |
  Analyze Go codebases for refactoring opportunities. Use when:
  (1) Reviewing code for shared functionality that can be factored out
  (2) Finding similar functions that can be combined using generics
  (3) Identifying opportunities for goroutines/channels to improve throughput
  (4) User asks to "review", "analyze", or "refactor" Go code
  (5) User wants to find code duplication or consolidation opportunities
---

# Go Code Review

Analyze Go code for refactoring opportunities: shared functionality extraction, generics consolidation, and concurrency improvements.

## Workflow

### Step 1: Run Linter with Auto-Fix

First, run golangci-lint to fix auto-fixable issues and identify remaining problems:

```bash
golangci-lint run --fix ./...
```

If lint errors remain after --fix, fix them manually before proceeding.

### Step 2: Analyze Target Code

Determine scope:
- If user specified files/packages: analyze those
- If no scope specified: ask user which packages or files to review
- For broad "review everything": focus on packages with recent changes or high complexity

Read the target files thoroughly before making recommendations.

### Step 3: Identify Refactoring Opportunities

Look for these patterns in priority order:

#### A. Shared Functionality Extraction

Find code that appears in multiple places with slight variations:
- Similar error handling blocks
- Repeated validation logic
- Common patterns

**Action**: Propose extracting to a shared function. Show before/after.

#### B. Generics Consolidation

Find functions/types that differ only in type parameters:
- Multiple functions doing the same thing for different types
- Slice utilities duplicated for different element types

**Action**: Propose generic function with type constraints. Prefer `[T any]` or `[T comparable]` over complex constraints.

#### C. Concurrency Improvements

Find opportunities for:
1. **Parallel independent operations**: Sequential loops that could run concurrently
2. **Fan-out/fan-in**: Multiple independent computations
3. **Context propagation**: Missing context.Context parameters or cancellation checks

### Step 4: Present Findings

For each opportunity found, provide:

1. **Location**: file:line
2. **Category**: Shared functionality / Generics / Concurrency
3. **Current code**: Brief snippet showing the pattern
4. **Proposed change**: Concrete refactored code
5. **Trade-offs**: Any downsides or considerations

### Step 5: Confirm Before Creating Packages

**CRITICAL**: If a refactoring would require creating a new package or folder, ASK THE USER FIRST.

## Code Style Requirements

When proposing changes, follow these project conventions:

- **External test packages**: Use `package foo_test` for tests
- **Strong typing**: Use enums/string constants, not magic values
- **Error wrapping**: Always wrap with context: `fmt.Errorf("operation: %w", err)`
- **Structured logging**: Use `log/slog` with structured fields
- **Minimal interfaces**: Only when genuinely needed for abstraction
- **No complex generics**: Prefer simple `[T any]` or `[T comparable]`

## What NOT to Do

- Don't add unnecessary abstractions
- Don't refactor code that isn't broken
- Don't add features beyond what's requested
- Don't create documentation files
- Don't add comments to unchanged code
- Don't create new packages without asking
