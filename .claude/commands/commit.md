---
description: Stage all changes and commit with a descriptive message
---

## Sync Beans to ClickUp

Before staging, sync beans in the background (non-blocking):

```bash
beanup --config .beans.clickup.yml sync &
```

## Stage and Commit

1. Run `git status --short` to see changes
2. Run `git diff HEAD` to review all changes
3. Stage all relevant changes
4. Commit with a concise, descriptive message:
   - Lowercase, imperative mood (e.g., "add feature" not "Added feature")
   - Focus on "why" not just "what"
   - Include affected bean IDs if applicable
5. Run `git status` to confirm the commit succeeded
