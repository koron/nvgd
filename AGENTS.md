# NVGD — Night Vision Goggles Daemon

Go HTTP file server for DevOps. Serves data from files, databases, S3, Redis, etc., with filters and transformations.

## Build & run

```
make build          # go build -gcflags '-e' ./...
make test           # go test ./...
make race           # go test -race ./...
make checkall       # go vet && staticcheck ./...
make vet            # go vet ./...
make staticcheck    # staticcheck ./...
make cover          # coverage report
make main-build     # build all main packages
```

**Requirements:** Go 1.25.0+, CGO enabled.

The `nvgd` binary is pre-built and in PATH. Run `nvgd` to start the server.

## Server

- Default address: `127.0.0.1:9280`
- Config file: `nvgd.conf.yml` (or `-c` flag)
- `nvgd -version` shows version

## Docker services

```
docker-compose up postgres mysql redis
# postgres: psql -U postgres
# mysql:   mysql -u mysql -p=abcd1234
# redis:   redis-cli
```

## Playwright tests

Tests live in `tests/`. Run with:

```
npx playwright test
```

No webServer config in `playwright.config.ts` — tests expect a running nvgd server. The Playwright MCP server (`npx playwright run-test-mcp-server`) is configured in `opencode.json`.
