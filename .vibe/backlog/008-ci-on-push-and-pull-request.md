---
status: todo
---
# CI on push and pull request

## Description
Run build, `go test ./...`, `go vet ./...`, `gofmt -l .` and the WASM build on every push and pull request, like the other Go library repos' CI items.

## Acceptance Criteria
- [ ] A failing test or lint fails the workflow on push and on pull requests
