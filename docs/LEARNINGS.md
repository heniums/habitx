# LEARNINGS — lessons and concept mastery

Two things live here.

1. **Concept ladder** — every concept the project teaches, with a level. Levels: `heard` (read about it) -> `help` (built it with AI help) -> `solo` (built it alone) -> `teach` (can explain it back without notes). A concept counts as **learned** only at `teach`. Tick the box when you reach `teach`.
2. **Lessons** — append anything that cost more than 15 minutes: date, area, one-line lesson. This list becomes the interview story and the personal docs.

## Concept ladder

### Go

- [ ] modules and toolchain (go run / test / fmt / vet) — level: solo
- [ ] structs, methods, interfaces — level: not started
- [ ] error values, wrapping with %w, errors.Is/As — level: not started
- [ ] net/http ServeMux routing and handlers — level: solo
- [ ] context.Context conventions — level: help
- [ ] database/sql + pgx — level: help
- [ ] table-driven tests, httptest — level: solo
- [ ] middleware chaining — level: not started
- [ ] JWT + bcrypt in Go — level: not started
- [ ] time and timezone handling — level: not started

### SQL

- [ ] DDL and migration discipline — level: solo
- [ ] joins and aggregates — level: not started
- [ ] composite PKs, idempotent writes — level: not started
- [ ] window functions (Phase 5) — level: not started

### React Native

- [ ] Expo + expo-router — level: not started
- [ ] View/Text/Pressable primitives — level: not started
- [ ] NativeWind styling — level: not started
- [ ] TanStack Query on mobile — level: not started
- [ ] FlatList and long lists — level: not started
- [ ] optimistic updates — level: not started
- [ ] secure token storage — level: not started
- [ ] dev build and release APK — level: not started

### Ops

- [ ] multi-stage Go Docker build — level: not started
- [ ] Caddy/Nginx TLS + same-origin proxy — level: not started
- [ ] CI gates (go test, golangci-lint, tsc) — level: not started

## Lessons

Format: `date - area - lesson`. Newest first.

- 2026-11-26 - go - struct tags take their quotes inside the backticks (`json:"name"`); `go vet` catches malformed tag grammar
- 2026-11-26 - http - POST creation answers 201 Created, not the default 200; the status code is part of the API contract
- 2026-10-01 - testing - Go compiles per directory; test files live in the same directory (and usually the same package) as the code they test
- 2026-10-01 - error handling - `if err := <some-function call>; err != nil {}` is a very common way quickly check error
- 2026-10-01 - git - gitignore bare-name matching at every depth; anchor with `/`; debug via `git check-ignore -v`
- 2026-09-22 - process - tracking lives in repo files (this set), not in any single AI tool's memory, so any helper can resume cold.
