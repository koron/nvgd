# AGENTS.md

NVGD ("Night Vision Goggles Daemon"): a single-binary Go HTTP server that serves
files, command output, S3, databases, and Redis through URLs like
`/file:///path?tail&grep=...`, with pluggable protocols and filters.

## Layout

- `main.go` → `config.LoadConfig` → `core.Run`. `core/` is the HTTP server,
  URL dispatch, alias, and filter pipeline.
- `protocol/<name>/` — one package per URL protocol (`file://`, `db://`, ...),
  each self-registers in `init()` via `protocol.Register`.
- `filter/<name>/` — one package per filter, self-registers via `filter.Register`.
- `plugins/` — only blank imports that wire default protocols/filters into
  the binary; `main.go` imports it.
- `internal/` — shared plumbing; `internal/protocoltest` and
  `internal/filtertest` are the standard test helpers.
- `README.md` is embedded (`go:embed`) and served at `/help/`, so README
  changes affect runtime behavior.

## Build & test

- `go.mod` requires Go 1.26 (README's "Go 1.24" is stale).
- **CGO is required** (`mattn/go-sqlite3`); tests won't build with `CGO_ENABLED=0`.
- `make test` / `go test ./...` — full suite. Single package:
  `make test TEST_PACKAGE=./filter/grep/` or
  `go test ./filter/grep/ -run TestX`.
- `make race`, `make cover`; `make checkall` = `go vet` + `staticcheck`
  (staticcheck uses `staticcheck.conf`, checks = all).
- Tests are self-contained: DB tests use in-memory sqlite, Redis tests use
  `miniredis`. No Docker services are needed for `go test` — `docker-compose.yml`
  (postgres/mysql/redis) is for manual protocol testing only.

## Running

- `go run .` — no config file needed; defaults come from `config/const_*.go`,
  which are **build-tag split**: without tags (release) default addr is
  `0.0.0.0:9280` with access log discarded; with `-tags debug` it is
  `127.0.0.1:9280` with access log on stdout.
- Local config is `nvgd.conf.yml` in CWD (or `-c FILE`); it is gitignored and
  not committed.
- Embedded HTML/JS assets (`protocol/*/assets/`, `core/assets/`) are
  `go:embed`ed, but each protocol wraps them with `internal/devfs`:
  `go run . -devfs.root <dir>` serves real files from
  `<dir>/<protocol>/assets/...` instead, so asset edits don't need a rebuild.

## Conventions

- New protocol: package under `protocol/`, register in `init()`, add the blank
  import to `plugins/protocols.go` (and `config.RegisterProtocol` for
  configurable ones). Tests must include
  `protocoltest.CheckRegistered`; `internal/protocoltest.Open/ReadAllString`
  is the standard harness (same pattern for filters).
- CI (`.github/workflows/go.yml`): on every push runs `go test ./...` and
  multi-platform CGO builds; tags `vX.Y.*` produce releases.
- `tmp/` and `_build/` are gitignored local scratch space.
