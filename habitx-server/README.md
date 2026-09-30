# habitx server

Go API, stdlib `net/http` only. Layout:

- `cmd/habitx-server/` — main.go, the only place that knows about main
- `internal/habits/` — habit handlers + store interface
- `internal/auth/` — users, tokens, middleware (Phase 3)
- `internal/server/` — mux wiring, middleware chain
- `db/migrations/` — hand-written SQL, applied by goose
- `db/queries/` — hand-written SQL, read by sqlc to generate Go

## Commands

- `make run` — start the server
- `make test` — run all tests
- `make tidy` — go mod tidy

sqlc + goose + docker-compose arrive in Phase 1.
