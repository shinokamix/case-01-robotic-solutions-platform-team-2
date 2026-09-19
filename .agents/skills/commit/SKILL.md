---
name: commit
description: Use when creating or reviewing a Git commit.
---

# Create a commit

1. Inspect the staged diff. Split changes that do not form one coherent result.
2. Write a lowercase imperative subject that states the result without relying
   on the diff or ticket.
3. Name the affected behavior, screen, service, command, or file. Add a
   component or condition when it distinguishes the change.
4. Use product and code terms, lowercased to fit the format. Choose specific
   verbs and describe the result, unless the implementation itself is the
   purpose.
5. Replace vague phrases such as `fix bug`, `update code`, `make changes`,
   `cleanup`, and `wip` with the concrete result.
6. Keep the subject under 72 characters. Omit the final period and Conventional
   Commits prefixes such as `feat:`, `fix:`, and `chore:`.
7. Add a lowercase body only when the reason or a non-obvious consequence does
   not fit in the subject.
8. Commit only when a reader can tell what changed and what it affects.

Examples:

- `prevent duplicate session refresh`
- `show empty state for missing search results`
