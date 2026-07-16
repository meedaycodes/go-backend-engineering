# ADR: Split Database Connection Pool into Separate Read and Write Pools

**Status**: Accepted  
**Date**: 2026-07-12

## Context

The application previously used a single `pgxpool.Pool` for all database operations — reads and writes shared the same connection pool. As request volume grows, read-heavy workloads (listing users, fetching profiles) compete with writes for the same limited connections, creating a bottleneck. Read replicas can absorb this read traffic, but only if the application routes reads and writes to separate pools.

## Decision

Split the single connection pool into two: a **write pool** targeting the primary database, and a **read pool** targeting the read replica. The `ReadWriteUserRepository` encapsulates this routing — mutating operations (`Save`, `Update`, `Delete`) use the write pool; query operations (`FindByID`, `FindAll`, `FindByEmail`) use the read pool. The routing decision is made at the repository method level, invisible to callers.

## Consequences

**Gains**
- Read traffic is offloaded to the replica, reducing load on the primary
- Each pool can be sized and tuned independently for its workload
- The application can scale read capacity horizontally by adding replicas without changing application code

**Trade-offs**
- Replication lag means a read immediately following a write may return stale data (read-your-own-writes anomaly)
- Two pools to configure, monitor, and keep healthy instead of one
- The current implementation points both pools at the same URL — actual replica routing requires infrastructure changes (a load balancer, PgBouncer, or RDS read endpoint)

## Future Considerations

- Replace direct pool connections with PgBouncer in transaction pooling mode to reduce connection overhead at scale
- Use a managed read endpoint (e.g. AWS RDS reader endpoint) as the read pool target so replica failover is handled automatically
- For write-after-read consistency requirements, route those specific reads to the write pool

---

# ADR: Security Hardening — TLS, RBAC, and Vault-Backed Secret Management

**Status**: Accepted  
**Date**: 2026-07-16

## Context

The service previously communicated over plain HTTP, used a flat JWT with no role information, and read the database password directly from a `.env` file. These three gaps represent distinct threat vectors:

1. **No TLS** — traffic between clients and the server is unencrypted, making it trivially interceptable on any shared network.
2. **No RBAC** — every authenticated user has identical access regardless of their role. There is no way to restrict sensitive operations (e.g. delete) to privileged users.
3. **Plaintext credentials in environment files** — `.env` files get committed accidentally, appear in CI logs, and are readable by anyone with filesystem access. A leaked `.env` means a leaked database.

## Decisions

### 1. TLS via `ListenAndServeTLS`

The HTTP server was changed from `ListenAndServe` to `ListenAndServeTLS`, loading a certificate and private key from `certs/cert.pem` and `certs/key.pem`. In development, a self-signed certificate is used. In production, these would be replaced by certificates issued by a trusted CA (e.g. Let's Encrypt via cert-manager on Kubernetes).

### 2. Role-Based Access Control

Two changes were made:
- `model.User` gained a `Role string` field, backed by a database migration (`000002_add_role_to_users`).
- The JWT `Auth` middleware was updated to extract the `role` claim and store it in the request context under `RoleKey`.
- A new `Authorize(requiredRole string)` middleware reads the role from context and returns `403 Forbidden` if it does not match. It is designed to be chained after `Auth`, not used standalone.

The `role` claim is set at signup time and embedded in the JWT — the middleware does not make a database call per request, keeping the hot path fast.

### 3. HashiCorp Vault for Secret Management

The database password is no longer read from `DATABASE_URL` in the environment. At startup, `config.Load()` calls `FetchDBPasswordFromVault`, which:
- Creates a Vault client pointing at `VAULT_ADDR` (default: `http://localhost:8200`)
- Authenticates with `VAULT_TOKEN`
- Reads the `password` field from the KV v2 secret at `secret/database`
- Uses it to construct the database connection string at runtime

In development, Vault runs as a Docker container in dev mode with a static root token. In production, the token would be replaced by a short-lived dynamic credential issued via Kubernetes auth or AWS IAM auth — the application code does not change.

## Consequences

**Gains**
- All traffic is encrypted in transit — TLS eliminates passive eavesdropping
- Role checks are enforced at the middleware layer, not scattered across handlers — a single place to audit
- The database password never appears in any file on disk, CI log, or environment variable export
- Vault provides a full audit log of every secret read — who accessed what and when

**Trade-offs**
- TLS adds a small latency overhead and requires certificate lifecycle management
- Vault is a new infrastructure dependency — if Vault is unreachable at startup, the service fails to start (fail-fast is intentional)
- Dev mode Vault is not durable — secrets are lost on container restart and must be re-seeded
- The current `role` claim is embedded in the JWT at issue time — changing a user's role requires re-issuing their token; the old token remains valid until expiry

## Future Considerations

- Replace the static dev token with Vault Agent running as a sidecar, which handles token renewal automatically
- Use Vault's dynamic database secrets engine to issue short-lived, per-service PostgreSQL credentials instead of a shared static password
- Add role invalidation by maintaining a token blocklist in Redis for immediate revocation when a role changes
- Replace self-signed certs with cert-manager and Let's Encrypt for automatic certificate rotation
