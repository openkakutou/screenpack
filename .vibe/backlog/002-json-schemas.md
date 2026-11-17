---
status: todo
depends_on: [001]
---
# Publish the JSON Schemas

## Description
Roadmap decision `029`: one JSON Schema for `screenpack.json` and one per screen, matching the spec, published in the repo at stable URLs so editors offer completion and validation.

## Acceptance Criteria
- [ ] Every example in the spec validates against the schemas
- [ ] A wrong type, unknown property or missing required field is rejected with a clear path
- [ ] Schemas carry a format version
