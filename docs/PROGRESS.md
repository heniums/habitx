# PROGRESS — session log

Append-only journal, newest entry first. Any agent (or human) finishing a session appends one entry using the template below. Readers: the newest 3 entries are the fastest way to resume context.

Template:

```
## YYYY-MM-DD — short title (who)
**Shipped:** files/features/decisions that landed
**Learned:** what the owner now understands that he did not before
**Next:** the single most important next action
```

## 2026-09-22 — onboarding (AI + owner)

**Shipped:** project charter. Decisions: Go API + Postgres from day 1 - stdlib net/http only - sqlc with owner hand-writing all SQL (deliberate break of the AI-authored-data-layer habit) - Expo + expo-router + TanStack Query + NativeWind - frozen thin v1 (CRUD, check-off, streaks) - JWT auth from day 1 (sequenced as Phase 3) - Expo Go dev loop, release APK sideload at v1 - $0 deploy target Oracle Cloud Always Free - future web UI via react-native-web, no stack change. Created AGENTS.md, docs/ROADMAP.md, docs/PROGRESS.md, docs/LEARNINGS.md.

**Learned:** baseline recorded — owner is strong in React/TS/Node/Express/Postgres ops; zero Go and React Native; self-identified weakest area is the data layer (has relied on AI-generated ORM/migration code without understanding it).

**Next:** scaffold Phase 0 on request — Go module with /healthz + hardcoded /habits + httptest, Expo app fetching it, make dev for both.
