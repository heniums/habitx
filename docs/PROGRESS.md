# PROGRESS — session log

Append-only journal, newest entry first. Any agent (or human) finishing a session appends one entry using the template below. Readers: the newest 3 entries are the fastest way to resume context.

Template:

```
## YYYY-MM-DD — short title (who)
**Shipped:** files/features/decisions that landed
**Learned:** what the owner now understands that he did not before
**Next:** the single most important next action
```

## 2026-11-26 — Phase 1 server-side complete (owner + AI, catch-up entry)

**Shipped:** Server half of Phase 1, merged via PRs #1 and #2. docker-compose Postgres + goose wiring; owner hand-wrote `0001_habits.sql` (first hand-written migration); all five habit queries in `db/queries/habits.sql` feeding sqlc codegen; `internal/store` (pgxpool + sqlc; AI scaffold, owner integrated); full CRUD handlers — `GET|POST /habits`, `GET|PUT|DELETE /habits/{id}`, 201 on create; DB-backed integration tests via goose-as-library `TestDB` (owner-authored, isolated per-run database) + handler endpoint tests. This entry back-fills October→November sessions that were never logged.

**Learned:** struct tags take their quotes inside the backticks (`json:"name"`) — `go vet` catches the malformed grammar; POST creation answers 201 Created, not the default 200; goose can run as a library to give each test run its own database.

**Next:** mobile half of Phase 1 — TanStack Query hooks, create-habit form, FlatList on the habits tab. Then the accept test: create a habit on the phone, restart the server, habit survives.

## 2026-10-01 — Phase 0 shipped (owner + AI)

**Shipped:** Phase 0 acceptance passed. Owner hand-wrote `main.go` (healthz + listHabits, method-pinned routes, `newMux()` shape), fixed 4 compiler errors, wrote table-driven tests (`TestRoutes` with 405 guard + `TestListHabitsBody` decode assertions) — all green uncached. Expo app rendered on phone via Expo Go in tunnel mode (LAN spinner = Wi-Fi isolation; `--tunnel` fixed). Fixed my own scaffold bug: unanchored `habitx-server` gitignore rule was silently un-tracking `cmd/habitx-server/` source — owner caught it.

**Learned:** directory=package (tests live next to code), `err :=` idiom inside `if`, gitignore anchoring + `git check-ignore -v`, `go build` output-name collision, `Id`→`ID` initialism convention. Ladder: modules/toolchain, net/http routing, table-driven tests → `solo`.

**Next:** Phase 1 — docker-compose Postgres, then owner writes `0001_habits.sql` (first hand-written migration) + goose wiring.

## 2026-09-30 — Phase 0 scaffolding (AI + owner)

**Shipped:** Go module scaffold (`go.mod`, `cmd/habitx-server`, `internal/{habits,auth,server}`, `db/{migrations,queries}`, Makefile, README) and Expo app (`habitx-mobile`, SDK 57, expo-router tabs template, demo tab and modal trimmed). Fixed a nested-git accident: `habitx-server/` had its own `.git` (reflog shows an initial commit with the scaffold files and a misleading "project charter" subject) — removed it and folded everything into the root monorepo. Toolchain verified: Go 1.27.1, Node 26.7.0, npm 11.19.

**Learned:** repo now teaches the monorepo rule — nested `.git` directories silently exclude files from the parent's `git add -A`.

**Next:** owner writes `cmd/habitx-server/main.go` (healthz + hardcoded habits handlers) and one httptest; AI reviews after.

## 2026-09-22 — onboarding (AI + owner)

**Shipped:** project charter. Decisions: Go API + Postgres from day 1 - stdlib net/http only - sqlc with owner hand-writing all SQL (deliberate break of the AI-authored-data-layer habit) - Expo + expo-router + TanStack Query + NativeWind - frozen thin v1 (CRUD, check-off, streaks) - JWT auth from day 1 (sequenced as Phase 3) - Expo Go dev loop, release APK sideload at v1 - $0 deploy target Oracle Cloud Always Free - future web UI via react-native-web, no stack change. Created AGENTS.md, docs/ROADMAP.md, docs/PROGRESS.md, docs/LEARNINGS.md.

**Learned:** baseline recorded — owner is strong in React/TS/Node/Express/Postgres ops; zero Go and React Native; self-identified weakest area is the data layer (has relied on AI-generated ORM/migration code without understanding it).

**Next:** scaffold Phase 0 on request — Go module with /healthz + hardcoded /habits + httptest, Expo app fetching it, make dev for both.
