# 🧠 System Design Interview Cheat Sheet
> Senior Backend Engineer — Tooters Delivery Interview Prep

---

## 📋 INTERVIEW FRAMEWORK (Always Follow This Order)
1. **Clarify Scope** — what to include/exclude
2. **Functional Requirements** — what the system does
3. **Non-Functional Requirements** — how well it does it
4. **Capacity Estimation** — scale, storage, bandwidth
5. **High-Level Architecture** — services, databases, queues
6. **Deep Dive** — critical components in detail
7. **Trade-offs** — justify every decision

---

## 📐 SCALABILITY

### Horizontal vs Vertical Scaling
- **Vertical** — bigger machine (CPU, RAM). Simple but has a limit.
- **Horizontal** — more machines. Complex but unlimited scale.
- **Rule of thumb:** Always design for horizontal scaling at senior level.

### Load Balancing
- Distributes traffic across servers.
- Algorithms: Round Robin, Least Connections, IP Hash (sticky sessions).
- Examples: AWS ALB, Nginx, HAProxy.

### Auto Scaling
- Automatically add/remove instances based on CPU, memory, or custom metrics.
- Use with stateless services — session state must be in Redis, not in-memory.

### Consistent Hashing
- Used in distributed caches/DBs to distribute data evenly.
- When a node is added/removed, only `K/N` keys need to be remapped (not all).
- Virtual nodes (vnodes) improve distribution — each physical node maps to multiple virtual nodes on the ring.
- Used by: Cassandra, DynamoDB, Redis Cluster.

---

## 🌐 NETWORKING

### DNS
- Translates domain names to IP addresses.
- TTL controls how long records are cached.
- Types: A record (IPv4), CNAME (alias), MX (mail).

### CDN (Content Delivery Network)
- Caches static assets (images, JS, CSS) at edge locations globally.
- Reduces latency by serving from the nearest node.
- Examples: CloudFront, Cloudflare, Akamai.

### Reverse Proxy
- Sits in front of servers, forwards client requests.
- Benefits: SSL termination, load balancing, caching, DDoS protection.
- Example: Nginx as reverse proxy.

### API Gateway
- Single entry point for all clients.
- Handles: Authentication, Authorization, Rate Limiting, Routing, Logging.
- Examples: AWS API Gateway, Kong.

### API Versioning
- **URL versioning:** `/v1/orders`, `/v2/orders` — simple, explicit, easy to route.
- **Header versioning:** `Accept: application/vnd.api+json;version=2` — cleaner URLs, harder to test.
- Always version from day one; deprecate with sunset headers and migration windows.
- Never break existing clients — additive changes only within a version.

---

## 💾 DATA

### SQL vs NoSQL
| | SQL | NoSQL |
|---|---|---|
| Schema | Fixed | Flexible |
| Scaling | Vertical (mostly) | Horizontal |
| ACID | Yes | Depends |
| Use when | Relations, consistency | Scale, flexibility |
| Examples | MySQL, PostgreSQL | MongoDB, Cassandra, DynamoDB |

### Database Sharding
- Split data across multiple DB instances (shards).
- Shard key choice is critical — avoid hotspots.
- Types: Range-based, Hash-based, Directory-based.
- Challenge: Cross-shard queries and joins are expensive.

### Replication
- **Master-Slave:** Master handles writes, slaves handle reads.
- **Master-Master:** Both handle reads/writes (conflict resolution needed).
- Provides: High availability + read scalability.

### Database Indexing
- Speeds up read queries at the cost of write performance and storage.
- **B-Tree index:** Default, good for range queries.
- **Hash index:** Good for equality lookups.
- **Spatial index:** For geospatial queries (lat/lng). Used with `ST_Distance_Sphere`.
- Tip: Index columns used in WHERE, JOIN, ORDER BY clauses.

### CAP Theorem
- A distributed system can only guarantee **2 of 3**:
  - **C**onsistency — all nodes see the same data
  - **A**vailability — system always responds
  - **P**artition Tolerance — system works despite network failures
- P is always required in distributed systems, so choice is **CP vs AP**.
- **CP** (e.g., order placement): Never return stale/wrong data.
- **AP** (e.g., order status updates): Slight staleness is acceptable.
- Apply CAP **per component**, not as a blanket rule across the whole system.

### ACID vs BASE
| ACID | BASE |
|------|------|
| Atomicity | Basically Available |
| Consistency | Soft state |
| Isolation | Eventually consistent |
| Durability | |
| SQL DBs | NoSQL DBs |
| Order placement | Order status feeds |

### Pagination Patterns
- **Offset pagination:** `LIMIT 20 OFFSET 100` — simple but slow at high offsets; breaks when rows are inserted mid-page.
- **Cursor-based (keyset) pagination:** `WHERE id > last_seen_id LIMIT 20` — O(1) regardless of depth, stable under inserts/deletes.
- **Rule of thumb:** Always use cursor-based at scale. Offset pagination is a red flag in senior interviews.

---

## ⚡ CACHING

### Cache Strategies
- **Write-through:** Write to cache AND DB simultaneously. Strong consistency, higher write latency.
- **Write-back (write-behind):** Write to cache first, async to DB. Fast writes, risk of data loss.
- **Write-around:** Write directly to DB, bypass cache. Good for infrequent reads.
- **Cache-aside (Lazy loading):** App checks cache first, on miss loads from DB and populates cache. Most common.

### Cache Eviction Policies
- **LRU (Least Recently Used):** Evict the item not used for the longest time. Most common.
- **LFU (Least Frequently Used):** Evict the item accessed least often.
- **TTL (Time To Live):** Expire items after a set duration.

### Redis vs Memcached
| Redis | Memcached |
|-------|-----------|
| Rich data structures (lists, sets, sorted sets, geo) | Simple key-value only |
| Persistence | No persistence |
| Pub/Sub support | No |
| Single-threaded | Multi-threaded |
| **Use Redis** for most cases | Use Memcached for simple, high-throughput caching only |

### Redis Geo
- `GEOADD` — add location
- `GEOSEARCH` — find points within radius
- Perfect for finding nearby restaurants/drivers.

---

## 📡 COMMUNICATION

### REST vs GraphQL vs gRPC
| | REST | GraphQL | gRPC |
|---|---|---|---|
| Protocol | HTTP | HTTP | HTTP/2 |
| Format | JSON | JSON | Protobuf (binary) |
| Flexibility | Fixed endpoints | Client defines query | Fixed contracts |
| Performance | Medium | Medium | High |
| Use when | Public APIs | Complex client needs | Internal microservices |

### WebSockets
- Full-duplex, bidirectional communication over a single TCP connection.
- Good for: Chat, live collaboration, multiplayer games.
- Not ideal for: Order status (one-way) — use SSE instead.
- **Caveat:** If the client ever needs to send data back (e.g. cancellation acknowledgements), SSE won't work — WebSockets or a separate REST call is required.

### SSE (Server-Sent Events)
- Unidirectional: server → client only.
- Works over standard HTTP, auto-reconnects.
- Perfect for: Order status updates, live feeds, notifications.
- Limitation: Client cannot push data back over the same connection.

### Long Polling vs SSE vs WebSockets
| | Long Polling | SSE | WebSockets |
|---|---|---|---|
| Direction | Server → Client | Server → Client | Bidirectional |
| Protocol | HTTP | HTTP | WS |
| Reconnect | Manual | Auto | Manual |
| Use for | Legacy systems | Live feeds | Chat, games |

### Message Queues
- Decouple producers from consumers.
- Async processing — producer doesn't wait for consumer.
- **RabbitMQ:** Simple pub/sub, routing, fanout exchanges. Good for task queues.
- **Kafka:** High-throughput, event streaming, message replay. Good at scale.
- **SQS:** Managed, serverless, integrates with AWS ecosystem.

### Dead Letter Queues (DLQ)
- Messages that repeatedly fail processing are routed to a DLQ instead of being lost or blocking the queue.
- Allows manual inspection, alerting, and reprocessing.
- Always configure a DLQ alongside any production queue — omitting it is a reliability gap.

### Backpressure
- What happens when consumers can't keep up with producers?
- **Solutions:** Bounded queues (reject or block when full), consumer auto-scaling, load shedding (drop low-priority messages), rate limiting at the producer.
- Failing to handle backpressure causes unbounded queue growth → OOM crashes.
- Mention this whenever you discuss async pipelines at senior level.

### Pub/Sub Pattern
- Publisher emits events without knowing who consumes them.
- Multiple consumers can subscribe to the same event.
- Used in order management: Order placed → Notification Service + Restaurant Service both consume.

### Webhooks
- HTTP callbacks — server notifies another server when an event happens.
- Used for: Payment confirmations, third-party integrations.

---

## 🛡️ RELIABILITY

### Rate Limiting & Throttling
- Protect services from abuse and overload.
- Algorithms: Token Bucket, Leaky Bucket, Fixed Window, Sliding Window.
- Implement at API Gateway level.
- Store rate limit counters in Redis.

### Circuit Breaker
- States: **Closed** (normal) → **Open** (failing, stop requests) → **Half-Open** (test recovery).
- Prevents cascading failures when a downstream service is down.
- Libraries: Resilience4j, Hystrix.

### Bulkhead Pattern
- Isolate failures by partitioning resources (thread pools, connection pools) per downstream service.
- If one dependency exhausts its pool, others are unaffected.
- Named after watertight compartments in ships — one breach doesn't sink the whole vessel.
- Use alongside Circuit Breaker for defence-in-depth.

### Retry Logic & Exponential Backoff
- On failure, retry after: 1s → 2s → 4s → 8s (with jitter to avoid thundering herd).
- Always set a max retry limit.
- Combine with Circuit Breaker.

### Idempotency
- Same request applied multiple times = same result.
- Use an **idempotency key** (UUID) per request.
- Server stores processed keys — ignores duplicates.
- Critical for: Order placement, payments.

### Timeouts & Fallbacks
- Always set timeouts on external service calls.
- Define fallback behavior: return cached data, default response, or graceful error.

### SLOs, SLIs, and Error Budgets
- **SLI (Service Level Indicator):** The actual measured metric (e.g. p99 latency, error rate).
- **SLO (Service Level Objective):** The target for that metric (e.g. 99.9% of requests < 200ms).
- **Error Budget:** The allowed failure margin (100% − SLO). Spend it on risk; when exhausted, freeze risky deployments.
- Mentioning this in senior interviews signals you think about reliability as a product property, not just an ops concern.

---

## 🏗️ ARCHITECTURE PATTERNS

### Monolith vs Microservices
| Monolith | Microservices |
|----------|---------------|
| Simple to develop | Complex to manage |
| Hard to scale parts independently | Scale each service independently |
| Single deployment | Independent deployments |
| Good to start with | Good at scale |

### Strangler Fig Pattern
- Incrementally migrate a monolith to microservices by routing specific paths/features to new services while the monolith handles the rest.
- Never rewrite everything at once — too risky.
- Route via API Gateway or reverse proxy; retire monolith paths one by one.

### Event-Driven Architecture
- Services communicate via events (not direct calls).
- Loose coupling — services don't know about each other.
- Built on message queues or event streams (Kafka).

### CQRS (Command Query Responsibility Segregation)
- Separate **write model** (Commands) from **read model** (Queries).
- Write DB optimized for consistency; Read DB optimized for speed.
- Example: Order writes go to MySQL, order history reads from Elasticsearch.

### Saga Pattern
- Manages distributed transactions across microservices.
- **Choreography:** Services react to each other's events (decentralized).
- **Orchestration:** A central saga coordinator directs each step (centralized).
- Each step has a **compensating transaction** to undo on failure.

### Two-Phase Commit (2PC) vs Saga
- **2PC:** Coordinator asks all participants to prepare, then commit. Strongly consistent but **blocking** — if the coordinator crashes mid-flight, participants are locked. Not partition-tolerant.
- **Saga:** Eventual consistency via compensating transactions. Non-blocking, partition-tolerant, but harder to reason about.
- **Rule of thumb:** Use Saga in distributed microservices. 2PC is a useful foil to explain *why* you're using Saga.

### Transactional Outbox Pattern
- Solves the dual-write problem (DB write + queue publish).
- Write event to an **outbox table** in the same DB transaction as the data.
- A relay process reads the outbox and publishes to the queue.
- Guarantees at-least-once delivery.

### Service Discovery
- Services register themselves (Consul, Eureka).
- Others discover them dynamically instead of hardcoded IPs.
- Types: Client-side discovery, Server-side discovery.

---

## 🌍 MULTI-REGION & GEO-DISTRIBUTION

### Active-Active vs Active-Passive
- **Active-Passive:** One region handles traffic; the other is a warm standby. Simple failover, but standby capacity is wasted and failover takes time.
- **Active-Active:** Multiple regions serve traffic simultaneously. Better latency and utilisation, but requires conflict resolution for cross-region writes.

### Data Residency & Replication Lag
- Some data (user PII, payments) must stay in a specific region due to regulation (GDPR, etc.).
- Cross-region replication introduces lag — reads from a secondary region may be stale.
- Design read paths to tolerate this (AP behaviour) or route writes and reads to the same region (CP behaviour).

### Geographic Sharding
- Assign users/drivers to a region based on location.
- Most requests stay local; only cross-region traffic (e.g. a user travelling) needs routing logic.
- Use consistent hashing to assign geographic zones to regional clusters.

---

## 🗄️ STORAGE

### Object vs Block Storage
| Object Storage | Block Storage |
|----------------|---------------|
| Files, images, videos | Raw disk volumes |
| Flat namespace (no folders) | Used by DBs and VMs |
| S3, GCS | AWS EBS, SAN |
| Highly scalable | Low latency |

### Blob Storage
- Binary Large Objects — images, videos, documents.
- Use S3 or Azure Blob Storage.
- Always store a URL in your DB, not the raw file.

### Data Partitioning
- **Horizontal (Sharding):** Split rows across nodes by key.
- **Vertical:** Split columns — separate tables for different data groups.
- **Directory-based:** Lookup table maps keys to shards.

### Time-Series Databases
- Optimized for time-stamped data (metrics, IoT, location history).
- Examples: InfluxDB, TimescaleDB, Prometheus.
- Use for: Driver location history, order event logs, system metrics.

---

## 🚀 ADDITIONAL CONCEPTS (Commonly Asked)

### Distributed Locking
- Prevent race conditions in distributed systems.
- Use Redis `SET NX PX` (set if not exists with expiry).
- Example: Only one service processes an order at a time.

### Database Connection Pooling
- Reuse DB connections instead of creating new ones per request.
- Tools: PgBouncer (Postgres), HikariCP (Java).
- Critical at high scale — DB connections are expensive.

### Observability (Monitoring, Logging, Tracing)
- **Metrics:** Prometheus + Grafana
- **Logging:** ELK Stack (Elasticsearch, Logstash, Kibana)
- **Distributed Tracing:** Jaeger, Zipkin (trace requests across microservices)
- **SLOs / Error Budgets:** Define reliability targets, not just dashboards.
- Always mention this in senior-level interviews.

### Security Best Practices
- HTTPS everywhere (TLS).
- JWT for stateless authentication.
- Rate limiting to prevent DDoS.
- Input validation to prevent SQL injection.
- Principle of least privilege for service permissions.

---

## 🎯 DELIVERY SYSTEM SPECIFIC CONCEPTS

### Geospatial Queries
- Store locations as `POINT(lng, lat)` with SPATIAL INDEX.
- MySQL: `ST_Distance_Sphere`, `ST_Within`.
- Redis: `GEOADD`, `GEOSEARCH` for fast radius lookups.

### Real-time Location Tracking
- Driver sends GPS update every 3-5 seconds.
- Store latest location in Redis (fast reads).
- Persist to time-series DB for history.

### Driver Matching
- Find nearest available driver using geospatial query.
- Consider: Distance, driver rating, current load.
- Use consistent hashing to assign geographic zones to servers.

### Order State Machine
- Model the order lifecycle explicitly as a state machine — not just a status string.
- States: `PLACED → ACCEPTED → PREPARING → PICKED_UP → IN_TRANSIT → DELIVERED → COMPLETED`
- Also model failure paths: `CANCELLED`, `FAILED`, `REFUND_PENDING`.
- Benefits: Prevents illegal transitions, simplifies event publishing, makes bugs obvious.
- Each state transition publishes an event to the message queue — downstream services (notifications, billing, analytics) consume them.

### Surge Pricing / Demand Sensing
- Detect demand spikes in near-real-time using sliding window counters (Redis) per geographic zone.
- Compare current order rate against a rolling baseline (e.g. same time last week).
- Apply a pricing multiplier when supply (available drivers) is low relative to demand (open orders).
- Rate-limit multiplier changes to avoid jarring price swings within a short window.
- This is a strong design discussion area specific to delivery platforms — few candidates prepare for it.

---

### 💡 SENIOR ENGINEER TIPS

1. **Always ask clarifying questions** before designing.
2. **Define requirements** before jumping to architecture.
3. **Do capacity estimation** — it justifies your tech choices.
4. **Apply CAP theorem per component**, not as a blanket rule.
5. **Talk about trade-offs** — there's no perfect solution.
6. **Mention observability** (logging, metrics, tracing, SLOs) — juniors forget this.
7. **Think about failure scenarios** — what happens when X goes down?
8. **Start simple, then scale** — don't over-engineer from the start.
9. **Name real technologies** — Redis, Kafka, MySQL, not just "a cache" or "a queue".
10. **Draw clear boundaries** between services — each service owns its data.
11. **Use cursor-based pagination**, not offset, at scale.
12. **Model state explicitly** — a state machine beats an enum column with no transition rules.
13. **Mention DLQs and backpressure** whenever you discuss async pipelines.
14. **Know why you're NOT using 2PC** — it makes you sound like you've debugged distributed transactions before.

### 💡 MY FAV DRINK
1. tea is number one
2. green tea is also ok
3. like hot chocolate as well

