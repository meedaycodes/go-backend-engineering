# Go Backend Engineering — 30-Day Intensive

## What This Is

This is a structured 30-day self-study program to reach top 1% Go backend/systems engineering skill. ~3 hours/day, 90 hours total. Each day has a focused topic and a project to ship.

**Stack**: Go, PostgreSQL, Redis, Docker, Kubernetes, Linux, CI/CD, System Design

## How Sessions Work

Each day lives in its own folder:
```
week-1/day-01-setup/
week-2/day-14-integration-testing/
...
```

When a new session starts, ask the user which day they are on, read the relevant day's folder, and pick up from where they left off.

## Learning Philosophy

- **The user is learning** — do not write code for them unless explicitly asked
- Explain concepts first, ask what they want to tackle next
- Guide them to write the code themselves; provide hints, explain patterns, review what they write
- When they do ask for code, keep it focused and explain every non-obvious decision

## Current Progress

- Week 1 (Days 1-7): Go Fundamentals & Software Craftsmanship — completed
- Week 2 (Days 8-14): Production-Grade Backend — in progress
  - Day 08: Auth system (JWT, bcrypt, protected routes) — done
  - Day 09: Config, structured logging (slog) — done
  - Day 10: API design, validation, rate limiting — done
  - Day 11: Redis caching, background email worker — done
  - Day 12: gRPC & Protobuf — done
  - Day 13: SDLC, golangci-lint, Makefile — done
  - Day 14: Integration & E2E Testing — done
- Week 3 (Days 15-21): Docker, Linux & Deployment — not started
- Week 4 (Days 22-30): System Design & Capstone — not started

## Day 14 — Integration & E2E Testing

**Goal**: Write a full integration test suite using real dependencies in Docker containers.

**Key concepts to cover**:
1. `testcontainers-go` — spin up real Postgres + Redis containers inside test code
2. `TestMain(m *testing.M)` — Go's test lifecycle hook for global setup/teardown
3. Running migrations programmatically against the test container
4. Table-driven E2E tests hitting the full HTTP router via `httptest`

**The project**: `week-2/day-14-integration-testing/`
- The app is a REST API with signup/login (JWT auth) and protected user CRUD routes
- Postgres via `pgx/v5`, Redis via `go-redis/v9`, router via `chi`
- An `InMemoryUserRepository` exists but should NOT be used for integration tests — the point is real containers
- Test file lives at `internal/integration/auth_test.go` (currently empty)
- Migrations are in `migrations/000001_create_users_table.up.sql`

**Distinction from unit tests**:
- Unit tests → mock/fake dependencies, test one layer in isolation
- Integration tests → real Postgres container + real Redis container, test the full stack end to end
