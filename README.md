# Go Backend Engineering

A 30-day intensive program building production-grade Go backend systems from the ground up. Every day ships working code — no toy examples, no tutorials followed passively.

## Stack

Go · PostgreSQL · Redis · Docker · Kubernetes · gRPC · GitHub Actions · Linux

## What This Covers

### Week 1 — Go Fundamentals & Software Craftsmanship
| Day | Topic | Project |
|-----|-------|---------|
| 01 | Go toolchain, modules, types, control flow | Hello world HTTP service with proper module structure |
| 02 | Structs, interfaces, custom error types | Interface-driven design with `errors.Is`/`errors.As` |
| 03 | Goroutines, channels, race detector, context | Fan-out/fan-in pipeline processing CSV data |
| 04 | Clean architecture, dependency injection, repository pattern | Layered REST API skeleton |
| 05 | Unit tests, table-driven tests, mocks, benchmarks | 80%+ test coverage on Day 04 project |
| 06 | `net/http`, `chi` router, middleware chains, graceful shutdown | CRUD REST API with logging, auth, recovery middleware |
| 07 | PostgreSQL, `pgx`, connection pooling, migrations | PostgreSQL-backed REST API |

### Week 2 — Production-Grade Backend
| Day | Topic | Project |
|-----|-------|---------|
| 08 | JWT auth, bcrypt, protected routes, refresh tokens | Full auth system (signup/login/protected routes) |
| 09 | Config management, structured logging (`slog`), log levels | Config-driven app with structured JSON logs |
| 10 | API design, request validation, rate limiting, pagination | Validation middleware + rate limiter |
| 11 | Redis caching (cache-aside), background worker pattern | Async email worker with Redis-backed cache |
| 12 | Protocol Buffers, gRPC server/client, interceptors | gRPC microservice with REST gateway |
| 13 | `golangci-lint`, pre-commit hooks, Makefile automation | Full CI toolchain with lint/test/build/migrate targets |
| 14 | `testcontainers-go`, integration tests, E2E API testing | Full integration test suite (real Postgres + Redis in containers) |

### Week 3 — Docker, Linux & Deployment *(in progress)*
| Day | Topic | Project |
|-----|-------|---------|
| 15 | Linux filesystem, permissions, process management, networking | Deploy Go binary on Linux VM via SSH |
| 16 | Dockerfile, multi-stage builds, image optimization, `trivy` | Production-optimized Docker image (<20MB) |
| 17 | Docker Compose, multi-service orchestration, health checks | Full stack: app + Postgres + Redis + queue |
| 18 | GitHub Actions, CI/CD pipeline, secrets management | lint → test → build → push → deploy pipeline |
| 19 | Kubernetes: Pods, Deployments, Services, ConfigMaps | App deployed to local K8s cluster |
| 20 | HPA, health probes, resource limits, Helm charts | Helm chart with rolling update strategy |
| 21 | Prometheus metrics, Grafana dashboards, log aggregation | Monitoring stack with custom Go metrics |

### Week 4 — System Design & Capstone *(upcoming)*
| Day | Topic | Project |
|-----|-------|---------|
| 22 | CAP theorem, event sourcing, CQRS, Saga pattern | Design doc for an event-driven system |
| 23 | Read replicas, sharding, PgBouncer | Read/write splitting in Go |
| 24 | OWASP, TLS, mTLS, HashiCorp Vault, RBAC | Secure service with TLS + Vault |
| 25 | `pprof`, escape analysis, load testing with `k6` | Profile, optimize, and benchmark a hot path |
| 26 | Circuit breakers, retries with backoff, distributed tracing | Circuit breaker + retry with exponential backoff |
| 27-28 | Capstone: production-ready microservice | gRPC + REST · PostgreSQL + Redis · Docker · K8s · CI/CD · 70%+ test coverage |
| 29 | Chaos engineering, fault injection, graceful degradation | Break the capstone, watch it recover |
| 30 | System design practice, ADRs, retrospective | System design document + growth plan |

## Skills Demonstrated

- **Go**: interfaces, goroutines, channels, context, error handling, generics
- **Auth**: JWT, bcrypt, refresh tokens, middleware-based route protection
- **Databases**: PostgreSQL with `pgx/v5`, connection pooling, migrations, transactions
- **Caching**: Redis cache-aside pattern, TTL, cache invalidation
- **Testing**: unit tests, table-driven tests, integration tests with `testcontainers-go`, E2E API tests
- **gRPC**: Protocol Buffers, server/client, interceptors, REST gateway
- **Observability**: structured logging (`slog`), Prometheus metrics, Grafana dashboards
- **Infrastructure**: Docker multi-stage builds, Docker Compose, Kubernetes, Helm
- **CI/CD**: GitHub Actions pipelines, golangci-lint, automated testing
- **System Design**: distributed systems patterns, database scaling, security hardening

## Project Structure

```
week-1/
  day-01-setup/
  day-02-types/
  ...
week-2/
  day-08-auth/
  day-14-integration-testing/
  ...
week-3/
week-4/
```

Each day is a self-contained Go module with its own `go.mod`, `Makefile`, and README.
