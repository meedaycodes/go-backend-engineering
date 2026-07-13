# Order Processing System — Design Document

## 1. System Overview

The system is used by a consumer goods retail company to automate product replenishment.
It monitors stock levels and triggers purchase orders before products fall below a configured
threshold, preventing out-of-stock situations.

---

## 2. CAP Trade-off

**Choice**: CP (Consistency over Availability)

**Reason**: The system must operate on the true, accurate count of each product. A stale
stock count could cause the system to miss a threshold breach, resulting in an out-of-stock
situation — the exact problem the system exists to prevent.

**Implication**: During a network partition, the system will reject order operations until
the current stock count can be confirmed across nodes. A brief period of unavailability is
preferable to placing orders based on incorrect inventory data.

---

## 3. Event Sourcing Model

All state changes are recorded as immutable events. The current stock level and order status
are derived by replaying these events.

| Event | Description |
|-------|-------------|
| `ThresholdConfigured` | A minimum stock level has been set for a product |
| `StockLevelChecked` | The current stock level was evaluated against the threshold |
| `ReplenishmentOrderPlaced` | A purchase order was raised with a supplier |
| `OrderConfirmedBySupplier` | The supplier acknowledged and accepted the order |
| `StockReceived` | The ordered goods arrived and were added to inventory |

---

## 4. CQRS Split

### Command Side (writes — change state, produce events)

| Command | Event Produced |
|---------|---------------|
| `SetProductThreshold` | `ThresholdConfigured` |
| `PlaceReplenishmentOrder` | `ReplenishmentOrderPlaced` |
| `ConfirmDelivery` | `StockReceived` |

### Query Side (reads — built from events, optimised for fast lookups)

| Query | Read Model |
|-------|------------|
| What is the current stock level for a product? | `CurrentStockLevel` |
| What replenishment orders are pending? | `PendingOrders` |
| Who approved a given order? | `OrderAuditLog` |

---

## 5. Architecture Decision Record

**Title**: Use CP consistency model for inventory data

**Status**: Accepted

**Context**: The replenishment system must trigger orders precisely when stock falls below a
configured threshold. Inaccurate stock counts — even briefly — can cause missed triggers
and out-of-stock situations, directly impacting the business.

**Decision**: Adopt a CP (Consistency + Partition Tolerance) model. During a network
partition, the system will become temporarily unavailable rather than serve potentially
stale stock counts.

**Consequences**:
- **Gain**: Stock counts are always accurate. Orders are never placed or withheld based on
  stale data. The system behaves correctly even under concurrent writes.
- **Give up**: During a network partition, warehouse operators cannot place or query orders
  until the partition heals. A short operational outage is the trade-off for data correctness.
- **Mitigation**: Implement alerting so operations teams are notified immediately when a
  partition is detected, minimising the duration of unavailability.
