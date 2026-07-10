# Production Go Backend Patterns

A reference implementation covering authentication, caching, gRPC, integration testing, Kubernetes deployment, and system design — built with Go, PostgreSQL, Redis, Docker, and Kubernetes.

Every module ships working, production-grade code. No toy examples.

## Stack

Go · PostgreSQL · Redis · Docker · Kubernetes · gRPC · GitHub Actions · Linux

## Modules

### Foundations
| Module | Patterns Covered | Deliverable |
|--------|-----------------|-------------|
| Go Fundamentals | Modules, types, interfaces, error handling | HTTP service with proper module structure |
| Concurrency | Goroutines, channels, race detector, context cancellation | Fan-out/fan-in pipeline processing CSV data |
| Clean Architecture | Dependency injection, repository pattern, layered design | REST API skeleton with clear separation of concerns |
| Testing | Table-driven tests, mocks, benchmarks, 80%+ coverage | Full test suite with `testify` |
| HTTP & Middleware | `chi` router, middleware chains, graceful shutdown | CRUD REST API with logging, auth, recovery middleware |
| Database | PostgreSQL with `pgx/v5`, connection pooling, migrations | PostgreSQL-backed REST API |

### Production Backend
| Module | Patterns Covered | Deliverable |
|--------|-----------------|-------------|
| Authentication | JWT, bcrypt, refresh tokens, protected routes | Full auth system (signup/login/protected routes) |
| Observability | Config management, structured logging (`slog`), log levels | Config-driven app with structured JSON logs |
| API Design | Request validation, rate limiting, pagination | Validation middleware + token bucket rate limiter |
| Caching & Workers | Redis cache-aside, TTL, cache invalidation, background workers | Async email worker with Redis-backed cache |
| gRPC | Protocol Buffers, server/client, interceptors | gRPC microservice with REST gateway |
| Developer Tooling | `golangci-lint`, Makefile automation, CI pipelines | Full CI toolchain with lint/test/build/migrate targets |
| Integration Testing | `testcontainers-go`, E2E API tests, real dependencies | Full integration test suite (real Postgres + Redis in containers) |

### Deployment & Infrastructure *(in progress)*
| Module | Patterns Covered | Deliverable |
|--------|-----------------|-------------|
| Linux | Filesystem, permissions, process management, networking | Deploy Go binary on Linux VM via SSH |
| Docker | Multi-stage builds, image optimization, `trivy` security scanning | Production-optimized Docker image (<20MB) |
| Docker Compose | Multi-service orchestration, health checks, volumes | Full stack: app + Postgres + Redis |
| CI/CD | GitHub Actions, secrets management, automated deployments | lint → test → build → push → deploy pipeline |
| Kubernetes | Pods, Deployments, Services, ConfigMaps, Secrets | App deployed to local K8s cluster |
| Kubernetes Advanced | HPA, health probes, resource limits, Helm charts | Helm chart with rolling update strategy |
| Observability | Prometheus metrics, Grafana dashboards, log aggregation | Monitoring stack with custom Go metrics |

### System Design & Resilience *(upcoming)*
| Module | Patterns Covered | Deliverable |
|--------|-----------------|-------------|
| Distributed Systems | CAP theorem, event sourcing, CQRS, Saga pattern | Design doc for an event-driven system |
| Database Scaling | Read replicas, sharding, PgBouncer connection pooling | Read/write splitting in Go |
| Security | OWASP, TLS, mTLS, HashiCorp Vault, RBAC | Secure service with TLS + Vault |
| Performance | `pprof`, escape analysis, load testing with `k6` | Profile, optimize, and benchmark a hot path |
| Resilience | Circuit breakers, retries with exponential backoff, distributed tracing | Circuit breaker + retry implementation |
| Capstone | Full production microservice | gRPC + REST · PostgreSQL + Redis · Docker · K8s · CI/CD · 70%+ test coverage |
| Chaos Engineering | Fault injection, graceful degradation, load shedding | Break the capstone, watch it recover |

## Skills

- **Go**: interfaces, goroutines, channels, context, error handling
- **Auth**: JWT, bcrypt, refresh tokens, middleware-based route protection
- **Databases**: PostgreSQL with `pgx/v5`, connection pooling, migrations, transactions
- **Caching**: Redis cache-aside pattern, TTL, cache invalidation
- **Testing**: unit tests, table-driven tests, integration tests with `testcontainers-go`, E2E API tests
- **gRPC**: Protocol Buffers, server/client, interceptors, REST gateway
- **Observability**: structured logging (`slog`), Prometheus metrics, Grafana dashboards
- **Infrastructure**: Docker multi-stage builds, Docker Compose, Kubernetes, Helm
- **CI/CD**: GitHub Actions pipelines, `golangci-lint`, automated testing
- **System Design**: distributed systems patterns, database scaling, security hardening

## Structure

Each module is a self-contained Go project with its own `go.mod`, `Makefile`, and production-ready code.

```
week-1/   — Go fundamentals & clean architecture
week-2/   — Production backend patterns
week-3/   — Docker, Linux & deployment
week-4/   — System design & capstone
```
