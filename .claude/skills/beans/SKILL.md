---
name: beans
description: Issue tracking with beans CLI. Use when (1) creating/updating beans, (2) setting parent/blocking relationships, (3) querying beans with GraphQL, (4) debugging relationship issues.
---

# Beans Issue Tracker

Agentic-first issue tracker for managing tasks, features, bugs, and milestones.

## Critical: Use Full IDs for Relationships

**ALWAYS use full IDs (with prefix) when setting `parent` or `blocking` relationships:**

```bash
# CORRECT - use full ID
beans update rxjc --parent go-bigq-z1nd --blocking go-bigq-z1nd

# WRONG - short ID breaks relationship lookups
beans update rxjc --parent z1nd --blocking z1nd
```

## Bean Types

| Type | Description |
|------|-------------|
| `milestone` | Target release or checkpoint |
| `epic` | Thematic container for related work |
| `feature` | User-facing capability or enhancement |
| `bug` | Something broken that needs fixing |
| `task` | Concrete piece of work to complete |

## Statuses

| Status | Description |
|--------|-------------|
| `draft` | Needs refinement before work can begin |
| `todo` | Ready to be worked on |
| `in-progress` | Currently being worked on |
| `completed` | Finished successfully |
| `scrapped` | Will not be done |

## Common Commands

### Create a bean
```bash
beans create "Title" -t feature -s todo -d "Description..."
```

### Update status
```bash
beans update abc123 --status in-progress
```

### Set relationships
```bash
beans update abc123 --parent go-bigq-xyz789
beans update abc123 --blocking go-bigq-xyz789
beans update abc123 --remove-parent
```

## GraphQL Queries

### Get bean with relationships
```bash
beans query '{ bean(id: "abc123") {
  title status body
  parent { id title }
  children { id title status }
  blockedBy { id title }
  blocking { id title }
} }'
```

### Find actionable beans
```bash
beans query '{ beans(filter: {
  excludeStatus: ["completed", "scrapped", "draft"],
  isBlocked: false
}) { id title status type priority } }'
```

## Tips

- Use `beans list` to see current workflow context
- Use `--json` flag for machine-readable output
- Bean IDs have format `go-bigq-xxxx` (4 character nanoid)
- Files stored in `.beans/` directory
- Always commit bean files with related code changes
