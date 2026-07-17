# ADR: Performance Engineering — Profiling, Benchmarking, and Load Testing

**Status**: Accepted  
**Date**: 2026-07-17

## Context

Before this day, the service had no performance baseline. There was no way to answer questions like:
- How fast is `CreateUser`?
- Where does CPU time go under load?
- What is the throughput limit of the server?
- Which function is the bottleneck?

Without measurement, any "optimisation" is guesswork. A systematic performance engineering approach requires three things: benchmarks for micro-level measurement, a profiler for macro-level CPU/memory analysis, and a load testing tool to generate realistic pressure.

## Decisions

### 1. pprof on a Separate Debug Port

Go's `net/http/pprof` package is exposed on port `6060` via a dedicated plain HTTP server, separate from the main TLS server on `8081`. A blank import (`_ "net/http/pprof"`) registers the handlers on `http.DefaultServeMux` as a side effect.

The separate port is intentional — the debug endpoint must never be reachable from the public internet. In production it would be firewalled to internal networks only, or replaced by a push-based profiler (e.g. Pyroscope).

### 2. Benchmarks in the Service Layer

`BenchmarkCreateUser` lives in `internal/service/user_service_bench_test.go` using Go's built-in benchmark framework (`testing.B`). It uses `NewInMemoryUserRepository` and a nil cache to isolate the service layer from I/O. `b.ResetTimer()` excludes setup from the measurement.

Result: `556ns/op`, `511 B/op`, `4 allocs/op` — this is the baseline for the pure business logic path, excluding database and cache I/O.

### 3. k6 for Load Testing

k6 was chosen over alternatives (wrk, Apache Bench, Locust) because:
- Scripts are JavaScript — readable and version-controllable
- Built-in checks, thresholds, and metrics output
- Supports TLS with self-signed cert bypass (`--insecure-skip-tls-verify`)
- Outputs structured results that map directly to SLO metrics (p90, p95 latency)

The load test script (`load_test.js`) simulates 10 concurrent users each signing up a unique user and listing all users, with a 1-second sleep between iterations.

### 4. Rate Limiter Increased for Load Testing

The per-IP rate limit was increased from 10 req/s to 1000 req/s for load testing. All k6 traffic originates from `127.0.0.1`, so the original limit blocked everything. In production, the limit stays at 10 req/s and load tests run from distributed IPs or bypass the rate limiter via an internal endpoint.

## Findings

CPU profile under load revealed:

| Function | CPU Time | % of Total |
|---|---|---|
| `blowfish.encryptBlock` (bcrypt) | 10.23s | 84.27% |
| `blowfish.ExpandKey` (bcrypt) | 1.09s | 8.98% |
| `syscall.rawsyscalln` (I/O) | 0.22s | 1.81% |
| Everything else | 0.60s | 4.94% |

**Conclusion**: 93% of all CPU time is bcrypt. The database, JSON encoding, JWT signing, middleware chain, and routing are all negligible. The service is CPU-bound on password hashing, not I/O-bound.

**Request throughput**: ~18 req/s with 10 VUs at `avg=40ms` latency. This is expected — bcrypt at DefaultCost (cost 10) takes ~60-80ms per hash, which sets a hard ceiling on signup throughput.

## Consequences

**Gains**
- A repeatable performance baseline now exists — future changes can be measured against it
- The profiler identifies the real bottleneck (bcrypt) rather than leaving engineers to guess
- The benchmark gives a sub-microsecond measurement of pure service logic performance

**Trade-offs**
- pprof adds a small always-on overhead (negligible in practice)
- The debug port (`6060`) must be carefully firewalled in production
- bcrypt cannot be optimised without reducing security — it is the intended bottleneck

## Future Considerations

- Offload password hashing to the email worker goroutine for async signup (accept the request, hash in background, notify via email when ready)
- Add a `make load-test` Makefile target that runs k6 as part of the CI performance gate
- Integrate Pyroscope or similar continuous profiler in staging to catch regressions before production
- Set k6 thresholds (`thresholds: { http_req_duration: ['p(95)<200'] }`) to fail the load test if latency degrades
- Consider argon2id instead of bcrypt for new deployments — better resistance to GPU cracking and tunable memory cost
