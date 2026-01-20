# Benchmark Results — Embedded Schema (WASAText)

## Query: Conversations sorted by last_message_at (limit 20)

### BEFORE index
- Plan: COLLSCAN + SORT
- executionTimeMillis: 1
- totalDocsExamined: 80
- totalKeysExamined: 0

### AFTER index (last_message_at: -1)
- Plan: IXSCAN + FETCH + LIMIT
- executionTimeMillis: 1
- totalDocsExamined: 20
- totalKeysExamined: 20

## Observations
- Before indexing, MongoDB scanned the full collection (COLLSCAN) and performed an in-memory sort.
- After creating the index on `last_message_at`, MongoDB used an index scan (IXSCAN), reducing examined documents from 80 to 20.
- Even on a small dataset, the execution plan is clearly optimized; this difference becomes critical at scale (many conversations).
