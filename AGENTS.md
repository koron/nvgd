# NVGD Agent Guide

## Build & Test

```
make build          # go build -gcflags '-e' ./...
make test           # go test ./...
make race           # go test -race ./...
make vet            # go vet ./...
make staticcheck    # staticcheck ./... (uses staticcheck.conf: checks = ["all"])
make cover          # generates tmp/cover.html
```

## Architecture

- **Plugin registration**: protocols and filters self-register via `init()` in their packages. The `plugins/` package imports all subpackages with blank imports to trigger registration at startup.
- **Entry**: `main.go` → `config.LoadConfig()` → `core.Run()`
- **Protocol packages**: `protocol/*/` — each implements `protocol.Protocol` interface (`Open(*url.URL)`)
- **Filter packages**: `filter/*/` — each implements `filter.Factory` signature (`func(*resource.Resource, filter.Params) (*resource.Resource, error)`)
- **Config**: `config/` — uses `goccy/go-yaml`. Default config file is `nvgd.conf.yml` (gitignored). Default listen: `127.0.0.1:9280`.
- **Resource**: `resource/` — abstraction wrapping `io.ReadSeekCloser` with metadata (content type, skip filters, etc.)

## Adding a Filter or Protocol

New filters/protocols must register themselves in their package's `init()`:

```go
func init() {
    filter.MustRegister("myfilter", myFilterFactory)
}
```

Then add the import to `plugins/filters.go` or `plugins/protocols.go`.

## Testing

- Go tests: `go test ./...` — all test files follow `*_test.go` convention
- Playwright tests: `npx playwright test` — tests in `tests/` directory
- Playwright config: `playwright.config.ts` — runs against chromium, firefox, webkit. No webServer configured (tests need an external server).

## Docker Services

```
docker-compose up postgres mysql redis
```

- PostgreSQL: `postgres://postgres:abcd1234@127.0.0.1:5432/postgres`
- MySQL: `mysql:abcd1234@tcp(127.0.0.1:3306)/mysql?sql_mode=TRADITIONAL`
- Redis: `redis://127.0.0.1:6379/0`

## Gotchas

- Go 1.25.0 required (per `go.mod`)
- `config.SecretString` marshals as `"__SECRET__"` in YAML output (don't assume plaintext)
- Config loading falls back to defaults if `nvgd.conf.yml` is missing (no error)
- `filter/base.go` sets `MaxLineLen: 1MB` default — can be overridden via `config.filters._base_.max_line_len`
- Filters are applied in query-string order; chain with `&`
- Default filters are skipped for "small" resources (per resource metadata)
- `-devfs.root` flag overrides embedded resources for development
