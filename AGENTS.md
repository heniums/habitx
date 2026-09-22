# AGENTS.md — habitx working agreement

Read this file first, every session. habitx is a habit-tracking app built **to learn Go and React Native by shipping a real product the owner uses daily**. The owner is the sole user. He has 2+ years of React/TypeScript/Node/Express/Postgres experience and **no prior Go or React Native experience**.

## The deal (learning rules — override convenience)

1. **The owner authors; AI reviews.** Never write his fundamentals for him: SQL migrations, queries, Go handlers, streak logic, auth code. He writes first; you review, explain, critique after. Scaffolding, config, and plumbing MAY be AI-generated.
2. **Explain, don't just fix.** Every fix ships with a one-or-two-sentence reason the owner can learn from.
3. **Tests where cheap:** pure logic (streaks, dates, middleware) gets table-driven tests. Never test scaffolding.
4. **No silent scope growth.** v1 scope is frozen in docs/ROADMAP.md; new ideas go to the Phase 5+ ladder, not into v1.

## Stack (locked 2026-09-22)

| Layer | Choice |
|---|---|
| Server | Go, stdlib `net/http` only (no framework), `internal/` layout |
| Data | Postgres in Docker, hand-written SQL, goose migrations, sqlc codegen |
| Auth | JWT (`golang-jwt/v5`) + bcrypt, single-user registration |
| Mobile | Expo + expo-router + TypeScript + TanStack Query v5 + NativeWind |
| Tokens | expo-secure-store (native), localStorage (web) via `Platform.OS` switch |
| Web UI | Phase 5+: react-native-web first; escape hatch = separate Next.js app sharing a `habitx-core` API package |
| Deploy | Oracle Cloud Always Free ARM (fallbacks: GCP e2-micro, home PC + Tailscale), Docker, Caddy same-origin `/api` proxy |

## Repo map

- `habitx-server/` — Go API (`cmd/`, `internal/`, `db/migrations/`, `db/queries/`, `sqlc.yaml`)
- `habitx-mobile/` — Expo app (`app/` routes, `src/api/` data layer)
- `docs/ROADMAP.md` — phases, deliverables, concepts, acceptance criteria. **The state of the project.**
- `docs/PROGRESS.md` — append-only session log, newest first. **What happened, what's next.**
- `docs/LEARNINGS.md` — lessons + concept mastery ladder. **The state of the learner.**

## Session protocol (any agent, any tool, any human)

**Start:** read `docs/ROADMAP.md` (current phase + acceptance criteria), then the newest 3 entries of `docs/PROGRESS.md`, then the `docs/LEARNINGS.md` ladder. Work within the current phase; do not skip ahead until its acceptance line passes.

**End (before finishing, no exceptions):**

1. Append a dated entry to `docs/PROGRESS.md`: shipped / learned / next.
2. Tick or advance checkboxes in `docs/ROADMAP.md`; update the current-phase line.
3. Move any lesson that cost more than 15 minutes into `docs/LEARNINGS.md`.
4. Commit everything. When a phase's acceptance criteria pass, tag `v0.N`.
