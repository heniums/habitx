# ROADMAP — habitx

The state of the project. Tick boxes as work completes; a phase is **done** only when its Accept line passes. Do not start a phase before the previous one ships. Phase 5+ items are picked by daily-use need, in any order.

**Current phase: 0 — not started**

## v1 = Phases 0-4

Frozen scope: register/login - habit CRUD (name, color, icon, daily schedule only) - today check-off - streak counter - per-habit history list. Nothing else belongs in v1.

API surface: `POST /users` - `POST /tokens` - `GET /me` - `GET|POST /habits` - `PATCH|DELETE /habits/{id}` - `PUT|DELETE /checkins/{habitID}/{date}` - `GET /habits/{id}/checkins` - `GET /healthz`. Check-in uses PUT because checking off twice must equal once (idempotency by design).

### Phase 0 — End to end in a weekend

- [ ] Go module boots, `GET /healthz` returns JSON
- [ ] `GET /habits` returns a hardcoded slice (no DB yet)
- [ ] Expo app runs on the phone via Expo Go, renders the list from the server (LAN IP)
- [ ] One handler covered by a table-driven httptest

Concepts: Go modules and toolchain - structs/methods - slices/maps - error values - JSON encoding - ServeMux routing - expo-router file tree - View/Text/StyleSheet - fetch-to-localhost trap (LAN IP)

**Accept: the phone shows the server's data.**

### Phase 1 — Persistence

- [ ] docker-compose runs Postgres
- [ ] Migration 0001 (habits) hand-written, applied via goose
- [ ] sqlc wired; all queries hand-written in `db/queries/`
- [ ] Habits CRUD endpoints behind a small store interface
- [ ] Integration test against the docker Postgres
- [ ] Mobile: TanStack Query hooks, create-habit form, FlatList, tab layout

Concepts: SQL DDL - database/sql + pgx - sqlc codegen - context.Context everywhere - interfaces defined at the consumer - query invalidation

**Accept: a created habit survives a server restart.**

### Phase 2 — The daily loop

- [ ] `check_ins` table with composite PK (habit_id, date)
- [ ] Idempotent PUT/DELETE check-in endpoints
- [ ] Streak computed by a pure, table-tested Go function, timezone-correct
- [ ] Today screen with optimistic check-off
- [ ] History list per habit

Concepts: composite PK - idempotent writes - UTC vs phone-local dates - optimistic updates

**Accept: the owner uses the dev build daily for 7 straight days.**

### Phase 3 — Auth (owner's explicit choice, sequenced here)

- [ ] users table, register, login issues JWT
- [ ] Auth middleware; user ID injected via context
- [ ] All habit/checkin routes locked; 401 handled on mobile with logout flow; token in expo-secure-store
- [ ] Table-driven tests for the middleware

Concepts: bcrypt - JWT sign/verify - middleware chaining - errors.Is/As + %w wrapping

**Accept: unauthenticated request gets 401; login works on the phone.**

### Phase 4 — Ship v1

- [ ] Multi-stage Dockerfile (small Go image)
- [ ] VPS live: Oracle free tier, Caddy TLS, same-origin `/api` proxy
- [ ] GitHub Actions: go test + golangci-lint + tsc
- [ ] Release APK built and sideloaded
- [ ] Tag v1.0

Concepts: Go Docker builds - reverse-proxy TLS - CI gates

**Accept: habitx is installed on the phone and used daily, anywhere.**

## Phase 5+ ladder (pick by need, any order)

- [ ] Reminders — expo-notifications; forces first dev build (Expo Go cannot deliver Android notifications)
- [ ] Heatmap + stats — SQL aggregates/window functions, SVG in RN
- [ ] Weekly / X-per-week schedules — schema evolution, migration discipline
- [ ] Offline-first — TanStack Query persistence, optimistic queue, then real sync
- [ ] Web UI — react-native-web evaluation; escape hatch: Next.js + shared `habitx-core` package
- [ ] Home-screen widget — native modules
