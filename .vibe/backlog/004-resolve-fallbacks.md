---
status: todo
depends_on: [003, 007]
---
# Resolve per-screen fallbacks to the default screenpack

## Description
Roadmap decision `029`: a screen missing from a screenpack, or an element missing from a screen, falls back to the default screenpack's.

## Acceptance Criteria
- [ ] A screenpack with only `title.json` resolves every other screen from the default
- [ ] The resolved model says where each screen came from
