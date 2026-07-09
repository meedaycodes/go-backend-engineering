# Go Systems Engineering — 30-Day Intensive

> **Goal**: Become a top 1% Go backend/systems engineer
> **Timeline**: 30 days, ~3 hours/day (90 total hours)
> **Covers**: Go, PostgreSQL, Redis, Docker, Kubernetes, Linux, CI/CD, System Design

---

## Week 1: Go Fundamentals & Software Craftsmanship (Days 1-7)

### Day 01 — Go Setup & Toolchain
- Go workspace, modules, `go build/run/test/vet/fmt`
- Variables, types, control flow, functions
- **Project**: Hello world HTTP service with proper module structure

### Day 02 — Types, Interfaces & Error Handling
- Structs, methods, interfaces, embedding
- Custom error types, `errors.Is`/`errors.As`, wrapping
- **Project**: Interface-driven design exercises

### Day 03 — Concurrency
- Goroutines, channels, `select`, `sync.WaitGroup`, `sync.Mutex`
- Race detector, context cancellation
- **Project**: Fan-out/fan-in pipeline processing CSV data

### Day 04 — Project Structure & Clean Architecture
- `cmd/`, `internal/`, `pkg/` layout
- Dependency injection, repository pattern
- **Project**: Scaffold a layered REST API skeleton

### Day 05 — Testing
- Unit tests, table-driven tests, subtests
- Mocks (`testify`), test helpers, benchmarks
- **Project**: Achieve 80%+ coverage on Day 04's project

### Day 06 — HTTP Servers & Middleware
- `net/http`, `chi` router, middleware chains
- Graceful shutdown, request/response patterns
- **Project**: CRUD REST API with logging, auth, recovery middleware

### Day 07 — Database Fundamentals
- PostgreSQL, `pgx`/`sqlx`, connection pooling
- Migrations with `golang-migrate`, transactions
- **Project**: Integrate PostgreSQL into the REST API

---

## Week 2: Production-Grade Backend (Days 8-14)

### Day 08 — Authentication & Authorization
- JWT tokens, bcrypt password hashing
- Auth middleware, protected routes, refresh tokens
- **Project**: Full auth system (signup/login/protected routes)

### Day 09 — Configuration, Logging & Observability
- Config management (`viper`, env vars, `.env`)
- Structured logging (`slog`/`zap`), log levels
- OpenTelemetry tracing basics
- **Project**: Config-driven app with structured logs and traces

### Day 10 — API Design Best Practices
- OpenAPI/Swagger spec, request validation
- Pagination, filtering, sorting patterns
- Rate limiting middleware
- **Project**: Swagger docs + validation middleware

### Day 11 — Caching & Message Queues
- Redis: caching strategies, TTL, cache invalidation
- Message queues: RabbitMQ or NATS
- Background worker pattern
- **Project**: Async email worker consuming from a queue

### Day 12 — gRPC & Protobuf
- Protocol Buffers, service definitions
- gRPC server/client, interceptors, streaming
- gRPC-Gateway (REST + gRPC)
- **Project**: gRPC microservice with REST gateway

### Day 13 — SDLC & Developer Tooling
- Git branching strategy (trunk-based / GitFlow)
- `golangci-lint`, pre-commit hooks
- Makefiles, task automation
- **Project**: Makefile with `lint`, `test`, `build`, `migrate` targets

### Day 14 — Integration & E2E Testing
- `testcontainers-go` for real dependencies
- Integration test patterns, test fixtures
- E2E API testing
- **Project**: Full integration test suite (real DB + Redis in containers)

---

## Week 3: Docker, Linux & Deployment (Days 15-21)

### Day 15 — Linux Fundamentals
- Filesystem hierarchy, permissions, users/groups
- Process management, systemd services
- Networking: `ss`, `curl`, `iptables` basics
- **Project**: Deploy a Go binary on a Linux VM via SSH

### Day 16 — Docker Deep Dive
- Dockerfile: multi-stage builds, layer caching
- Image optimization (<20MB), security scanning (`trivy`)
- `.dockerignore`, non-root users
- **Project**: Production-optimized Dockerfile

### Day 17 — Docker Compose
- Multi-service orchestration
- Networking, volumes, health checks
- Environment-specific overrides
- **Project**: `docker-compose.yml` running full stack (app + DB + Redis + queue)

### Day 18 — CI/CD with GitHub Actions
- Workflow syntax, triggers, matrix builds
- Pipeline: lint -> test -> build -> push image -> deploy
- Secrets management, caching dependencies
- **Project**: Working CI/CD pipeline

### Day 19 — Kubernetes Fundamentals
- Pods, Deployments, Services, ConfigMaps, Secrets
- `kubectl` commands, namespaces
- `minikube` or `kind` local cluster
- **Project**: Deploy app to local K8s cluster

### Day 20 — Kubernetes Advanced
- Horizontal Pod Autoscaler (HPA)
- Health probes (liveness, readiness, startup)
- Resource limits/requests, Ingress
- Helm charts
- **Project**: Helm chart with rolling update strategy

### Day 21 — Kubernetes Observability
- Prometheus: metrics exposition from Go (`promhttp`)
- Grafana dashboards, alerting rules
- Log aggregation patterns
- **Project**: Monitoring stack with custom Go metrics

---

## Week 4: System Design & Capstone (Days 22-30)

### Day 22 — Distributed Systems Patterns
- CAP theorem, eventual consistency
- Event sourcing, CQRS
- Idempotency, distributed transactions (Saga)
- **Project**: Design doc for an event-driven system

### Day 23 — Database Scaling
- Read replicas, write-ahead log
- Sharding strategies, partitioning
- Connection pooling at scale (PgBouncer)
- **Project**: Implement read/write splitting in Go

### Day 24 — Security Hardening
- OWASP Go Top 10, input sanitization
- TLS termination, mTLS between services
- Secret management (HashiCorp Vault)
- RBAC implementation
- **Project**: Secure service with TLS + Vault

### Day 25 — Performance Engineering
- `pprof`: CPU, memory, goroutine profiling
- Memory optimization, escape analysis
- Load testing with `k6` or `vegeta`
- **Project**: Profile, optimize, and benchmark a hot path

### Day 26 — Microservices Resilience
- Service discovery, circuit breakers (`gobreaker`)
- Retries with exponential backoff + jitter
- Distributed tracing across services
- Bulkhead and timeout patterns
- **Project**: Implement circuit breaker + retry with backoff

### Day 27-28 — Capstone Project (Part 1 & 2)
Build a **production-ready microservice** (choose one):
- URL shortener with analytics
- Distributed task queue
- Notification service (email/SMS/push)

**Requirements**:
- Clean architecture with dependency injection
- gRPC + REST API (dual interface)
- PostgreSQL + Redis
- Dockerized with multi-stage build
- Kubernetes manifests or Helm chart
- CI/CD pipeline (GitHub Actions)
- Prometheus metrics + structured logging
- 70%+ test coverage (unit + integration)

### Day 29 — Chaos Engineering & Resilience Testing
- Fault injection in your capstone
- Graceful degradation under failure
- Load shedding, circuit breaking in action
- **Project**: Break your capstone, watch it recover

### Day 30 — System Design & Retrospective
- Practice system design: rate limiter, distributed cache, or pub/sub system
- Write Architecture Decision Records (ADRs)
- Review all 30 days, identify gaps
- **Project**: System design document + personal growth plan

---

## Principles to Follow Every Day

1. **Ship working code daily** — commit with meaningful messages
2. **Test-Driven Development** — write tests before or alongside code
3. **Read before you write** — study source code of libraries you use
4. **Document decisions** — write ADRs for non-obvious choices
5. **Review your own code** — read your PRs critically before considering them done

## Recommended Resources

| Type | Resource |
|------|----------|
| Book | *Let's Go Further* — Alex Edwards |
| Book | *Designing Data-Intensive Applications* — Martin Kleppmann |
| Book | *The Go Programming Language* — Donovan & Kernighan |
| Repo | [golang-standards/project-layout](https://github.com/golang-standards/project-layout) |
| Tool | [Go Playground](https://go.dev/play/) for quick experiments |
| Docs | [Effective Go](https://go.dev/doc/effective_go) |
