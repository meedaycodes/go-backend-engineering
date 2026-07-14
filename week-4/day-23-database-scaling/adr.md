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

