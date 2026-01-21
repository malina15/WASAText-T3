# Benchmarks & results

## Dataset setup (embedded schema)

### MongoDB container with mounted scripts
![MongoDB container with mounted scripts](../images/screenshots/03-mongo-container-mounted-scripts.png)

### Connected to database and collections
![Connected to wasatext_embedded](../images/screenshots/04-mongosh-connected-collections.png)

### Clean database and populate dataset
![Drop DB and populate](../images/screenshots/05-dropdb-populate-counts.png)

## Benchmark query

The following query retrieves the 20 most recently active conversations:

```js
db.conversations.find({}).sort({ last_message_at: -1 }).limit(20).explain("executionStats")
```

### Explain BEFORE index
![Explain BEFORE index](../images/screenshots/06-explain-before-index.png)

### Create index on last_message_at
![Create index](../images/screenshots/07-create-index-last_message_at.png)

### Explain AFTER index
![Explain AFTER index](../images/screenshots/08-explain-after-index.png)

## Results table

| Scenario | Schema | Index | Plan (winning) | totalDocsExamined | totalKeysExamined | executionTimeMillis | Notes |
|---|---|---:|---|---:|---:|---:|---|
| Most recent conversations (`sort last_message_at desc`, `limit 20`) | Embedded | No | `COLLSCAN` + blocking `SORT` | 80 | 0 | 2 | Sort performed in memory (`usedDisk: false`) |
| Most recent conversations (`sort last_message_at desc`, `limit 20`) | Embedded | Yes (`last_message_at_-1`) | `IXSCAN` → `FETCH` → `LIMIT` | 20 | 20 | 0 | Index provides order; avoids blocking sort |

## Interpretation

- **Before indexing**, MongoDB scans the entire `conversations` collection (`COLLSCAN`, `totalDocsExamined = 80`) and then performs a **blocking in-memory sort** on `last_message_at` to return the top 20 results.
- **After creating the descending index** on `last_message_at`, the query switches to an **index scan** (`IXSCAN`) followed by `FETCH` and `LIMIT`. Because the index already stores documents in the needed order, MongoDB can return the first 20 entries directly, reducing work from scanning **80 docs** to **20 docs**.
- This benchmark demonstrates a typical messaging “recency feed” optimization: indexing the sort key reduces scanned documents and removes the blocking sort stage, improving latency and scalability as the dataset grows.
- `totalDocsExamined` decreased from **80** to **20** (4× fewer documents scanned).
