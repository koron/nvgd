# AGENTS.md

## What this is

**NVGD** (Go HTTP file server) + **Playwright tests** (TypeScript).

## Quick commands

| Task | Command |
|------|---------|
| Build (Go) | `go build ./...` or `make build` |
| Test all (Go) | `go test ./...` or `make test` |
| Race test (Go) | `make race` |
| Coverage (Go) | `make cover` |
| Vet (Go) | `go vet ./...` or `make vet` |
| Staticcheck | `staticcheck ./...` or `make staticcheck` |
| Lint + vet | `make checkall` |
| Playwright tests | `npx playwright test` |
| Install Playwright browsers | `npx playwright install --with-deps` |
| Single Go test | `go test -run TestName ./package/path` |
| Bench (Go) | `make bench` |

## Repo structure

- `main.go` — app entrypoint
- `config/` — YAML config loading (`nvgd.conf.yml`, `-c` flag)
- `core/` — server, routing, resource handling
- `protocol/` — protocol handlers (file, s3, db, redis, command, duckdb, etc.)
- `filter/` — response filters (grep, tail, head, markdown, htmltable, etc.)
- `plugins/` — plugin registration (wires protocols + filters)
- `internal/` — shared test helpers (`filtertest`, `protocoltest`)
- `resource/` — resource type shared across protocol/filter
- `tests/` — Playwright test specs
- `doc/` — supplementary docs
- `specs/` — test plans (README only)
- `.github/workflows/` — `go.yml` (build+release) and `playwright.yml` (test)

## Go backend

- Requires **Go 1.24+** with **CGO enabled**.
- Config file defaults to `./nvgd.conf.yml`; override with `-c path`.
- Dev flag: `-devfs.root <dir>` serves embedded resources from filesystem.
- Pprof flag: `-pprofaddr :6060`.
- Staticcheck config in `staticcheck.conf`: checks = `["all"]`.
- All `main` packages are cross-compiled and released on tag push.

## Playwright tests

- Config: `playwright.config.ts` (chromium, firefox, webkit).
- CI runs on push/PR to main/master.
- `npm ci` then `npx playwright install --with-deps` then `npx playwright test`.
- Tests are currently **stub/skeleton** (`tests/example.spec.ts`, `tests/seed.spec.ts`).
- OpenCode config (`opencode.json`) provides Playwright subagents for test generation/healing/planning.

## Docker (for DB protocol tests)

- `docker-compose up` starts PostgreSQL (5432), MySQL (3306), Redis (6379).
- Needed to run DB protocol handler tests that require live databases.

## Style conventions

- `go vet ./...` + `staticcheck ./...` must pass (run `make checkall`).
- No comments in code per project convention.
- Avoid generated files; no snapshot workflows.
- Go tests use `testing` stdlib; Playwright tests use `@playwright/test` with `expect`.
