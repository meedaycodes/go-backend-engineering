# ADR: Microservices Resilience — Circuit Breaker and Retry with Exponential Backoff

**Status**: Accepted  
**Date**: 2026-07-19

## Context

When a downstream service fails or becomes slow, requests from the calling service pile up waiting for responses. Without a mechanism to stop this, threads exhaust, memory grows, and the failure cascades — a single unhealthy downstream can bring down the entire system. The service needed a way to detect downstream failure early and stop sending requests to a service that is known to be down.

## Decisions

### 1. Circuit Breaker from Scratch

The circuit breaker was implemented from scratch rather than using a library, to understand the state machine directly. The circuit has three states:

- **Closed** — requests flow through normally; failures are counted
- **Open** — `maxFailures` has been reached; all requests are immediately rejected without hitting the downstream
- **HalfOpen** — the recovery `timeout` has elapsed since the last failure; one request is allowed through to test if the downstream has recovered. On success the circuit closes; on failure it opens again

Config used: `maxFailures = 3`, `timeout = 10s`. The mutex is released before calling the downstream function and re-acquired after — holding the lock during I/O would block all concurrent callers.

### 2. Retry with Exponential Backoff and Jitter

Retry accepts a `maxAttempts` value and retries the function on failure up to that limit. The wait between attempts grows exponentially from a 100ms base: `wait = 100ms * 2^attempt`. Jitter multiplies the wait by a random factor between 0.5 and 1.5 so that clients spread their retries across time rather than all hitting the recovering service at the same millisecond (thundering herd).

### 3. Nesting Order: Circuit Breaker Wraps Retry

The initial implementation had retry wrapping the circuit breaker. This meant that when the circuit was open, retry would still sleep its full backoff between fast-fail attempts — a 4th call still took 281ms instead of being instant. 

Swapping the order so the circuit breaker wraps retry fixes this: the circuit breaker checks state first and fast-fails in microseconds if open. If the circuit is closed, retry handles transient failures within that single attempt. The 4th call after the circuit tripped completed in 6µs.

## Consequences

**Gains**
- Fast-fail protects the service under load — no threads wasted waiting on a known-dead downstream
- Retry handles transient blips transparently without surfacing errors to the caller
- Jitter prevents thundering herd when many clients retry simultaneously after a shared outage

**Trade-offs**
- The `/call-downstream` route is temporarily unavailable while the circuit is open, even if only some requests would have succeeded
- Retry adds latency before giving up — 3 attempts with backoff can take ~700ms before returning an error on a genuinely dead service

## Future Considerations

- Emit metrics on circuit state transitions (open/close events) so oncall can observe how often the circuit trips in production
- Use per-route circuit breakers instead of one shared instance — a slow payments service should not affect the user profile service
- Distinguish retryable errors from non-retryable ones — a 400 Bad Request should never be retried; only 5xx and network errors should trigger retry
