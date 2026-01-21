# Sharding (T3) — Plan and justification (not implemented)

This project focuses on schema design, indexing, aggregation, and replication with a realistic WASAText-like workload.
Sharding is included here as a **scaling strategy plan** with clear design choices and evaluation criteria.

**Status:** Not implemented (see constraints at the end).  
**Goal of this document:** show how sharding would be applied to this dataset and how it would be evaluated.

---

## 1. What would be sharded and why

### Target collection
`conversations`

Reason:
- It is the largest and most frequently accessed collection.
- It contains the embedded `messages` array and related subdocuments (reactions, receipts), which dominates storage size.
- Key access patterns (Inbox sort by `last_message_at`, conversation view by `_id`) are centered around this collection.

### Expected sharding benefits
- Horizontal scale-out for storage and read throughput.
- Distributing large conversation documents across multiple machines.
- Reduced single-node resource pressure (CPU / disk / memory) for read-heavy workloads.

---

## 2. Workload assumptions (WASAText-like)

The following operations guide the sharding design:

1) Inbox query (already benchmarked for indexing):
- Fetch recent conversations ordered by `last_message_at`.
- Example pattern: `find({ ... }).sort({ last_message_at: -1 }).limit(k)`

2) Conversation view:
- Fetch a single conversation by `_id` with its recent messages.

3) Message write path:
- Append a new message to one conversation document and update `last_message_at`.

---

## 3. Shard key options and trade-offs

Choosing the shard key is the central decision. The dataset schema suggests these candidates:

### Option A (simple and robust): hashed `_id`
**Shard key:** `{ _id: "hashed" }`

Pros:
- Very good uniform distribution (avoids hotspots).
- Simple, works immediately with the existing schema (no schema changes required).
- Most operations that target a conversation by `_id` become targeted to a single shard.

Cons:
- Inbox query sorted by `last_message_at` becomes a scatter-gather across shards.
- Global sort + limit typically requires merging results from all shards (unless extra design is added).

When it fits best:
- Workloads dominated by fetching a conversation by ID and appending messages.

---

### Option B (tenant-based): `tenant_id` (ideal, but requires schema change)
**Shard key:** `{ tenant_id: 1, _id: 1 }` (or `{ tenant_id: "hashed" }`)

Pros:
- Excellent for multi-tenant scaling (each tenant’s data is colocated).
- Many queries become targeted by tenant (no scatter-gather).

Cons:
- Requires adding a `tenant_id` field to conversations and enforcing it in all queries.
- Not available in the current dataset/app design.

When it fits best:
- Real production deployments with multiple organizations/tenants.

---

### Option C (time-based): `last_message_at`
**Shard key:** `{ last_message_at: 1 }` (or compound with `_id`)

Pros:
- Potentially helps Inbox queries if they are range-based on time (e.g., “last 7 days”).

Cons:
- High risk of write hotspots (new messages update `last_message_at` constantly).
- Chunk migrations can become frequent and expensive.
- In this dataset, every new message updates `last_message_at`, so writes concentrate on the “most recent” key space.

Conclusion:
- Not recommended for the WASAText write pattern.

---

## 4. Proposed design for this project

**Chosen shard key (proposal):** `{ _id: "hashed" }`

Reason:
- Works with the current schema without changes.
- Supports the most common targeted query: fetch a conversation by `_id`.
- Avoids write hotspots that would appear with a time-based shard key.

### Mitigation for Inbox query
Inbox is a global sort by `last_message_at` and would likely scatter across shards.

Mitigations (conceptual):
- Maintain a secondary “inbox index” materialized view per user/participant (separate collection, smaller documents).
- Or query per user/participant if the application supports it (requires filtering by participant user_id).
- Keep indexing on `last_message_at` on each shard for local efficiency.

---

## 5. How sharding would be deployed (conceptual)

A minimal sharded cluster typically needs:
- **Config servers:** 3 nodes (replica set)
- **Shard replica sets:** e.g., 2 shards × 3 nodes each (for redundancy)
- **mongos routers:** 1 or more

Docker-based layout (conceptual):
- configsvr: cfg1, cfg2, cfg3 (replSet: cfgRS)
- shard1: s1a, s1b, s1c (replSet: shard1RS)
- shard2: s2a, s2b, s2c (replSet: shard2RS)
- mongos: router service

Shard enabling steps (conceptual mongosh):
```js
sh.enableSharding("wasatext_embedded")
sh.shardCollection("wasatext_embedded.conversations", { _id: "hashed" })
```

---

## 6. Evaluation plan (what would be measured)

To satisfy “apply and evaluate”, the following would be collected:

### Functional checks
- `sh.status()` shows:
  - sharding enabled for the database
  - `conversations` sharded with the chosen shard key
  - multiple shards registered

### Distribution checks
- `db.conversations.getShardDistribution()` (or equivalent stats) to show data spread.

### Performance checks (compare sharded vs non-sharded)
- Conversation view by `_id`:
  - Should target a single shard (expected improvement at scale).
- Message append:
  - Check write latency and whether writes are evenly distributed.
- Inbox query:
  - Measure scatter-gather cost and compare mitigation strategies (if implemented).

---

## 7. Why sharding is not implemented here (constraints)

Sharding was not implemented in this project due to practical constraints:

1) **Complexity and time budget**
- A correct sharded cluster requires multiple replica sets, config servers, and mongos.
- The setup is significantly more complex than a replica set demo.

2) **Evaluation quality**
- With a small synthetic dataset (80 conversations), sharding overhead dominates and does not provide meaningful performance insights.
- A fair evaluation would require a much larger dataset and repeatable benchmarking under load.

3) **Project scope**
- The project already demonstrates applied scaling via indexing and replication, both with measurable effects and clear demos.
- Sharding is documented as a plan, including key design choices and an evaluation methodology, but not executed.

---

## 8. Summary

- Sharding target: `conversations`
- Proposed shard key: `{ _id: "hashed" }`
- Trade-off: best for targeted `_id` access; Inbox query becomes scatter-gather
- Clear deployment outline + evaluation plan included
- Not implemented due to setup complexity and limited dataset size for a fair evaluation
