# Aggregation pipelines (T3) — WASAText embedded schema

All pipelines were executed in `mongosh` on database `wasatext_embedded`.

---

## Pipeline 1 — Top 10 users by number of messages sent

This pipeline computes the most active senders across all conversations.

```js
db.conversations.aggregate([
  { $unwind: "$messages" },
  { $match: { "messages.sender.user_id": { $ne: null } } },
  { $group: { _id: "$messages.sender.user_id", messages_sent: { $sum: 1 } } },
  { $sort: { messages_sent: -1 } },
  { $limit: 10 }
])
```

Output:
```js
[
  { _id: ObjectId('6971252004739cc1d7284d4c'), messages_sent: 314 },
  { _id: ObjectId('6971252004739cc1d7284d44'), messages_sent: 288 },
  { _id: ObjectId('6971252004739cc1d7284d2b'), messages_sent: 257 },
  { _id: ObjectId('6971252004739cc1d7284d21'), messages_sent: 221 },
  { _id: ObjectId('6971252004739cc1d7284d0c'), messages_sent: 221 },
  { _id: ObjectId('6971252004739cc1d7284daf'), messages_sent: 210 },
  { _id: ObjectId('6971252004739cc1d7284d86'), messages_sent: 187 },
  { _id: ObjectId('6971252004739cc1d7284db5'), messages_sent: 183 },
  { _id: ObjectId('6971252004739cc1d7284d0e'), messages_sent: 176 },
  { _id: ObjectId('6971252004739cc1d7284d36'), messages_sent: 173 }
]
```

---

## Pipeline 2 — Top 10 conversations by total number of messages

The dataset includes a precomputed `message_count` field per conversation, so this
pipeline sorts by that value.

```js
db.conversations.aggregate([
  { $project: { participants_count: 1, message_count: 1, last_message_at: 1 } },
  { $sort: { message_count: -1 } },
  { $limit: 10 }
])
```

Output:
```js
[
  {
    _id: ObjectId('6971252104739cc1d728c12d'),
    participants_count: 2,
    last_message_at: ISODate('2026-01-19T06:51:36.212Z'),
    message_count: 249
  },
  {
    _id: ObjectId('6971252104739cc1d728c13e'),
    participants_count: 6,
    last_message_at: ISODate('2025-12-27T08:43:26.831Z'),
    message_count: 247
  },
  {
    _id: ObjectId('6971252104739cc1d728c13a'),
    participants_count: 3,
    last_message_at: ISODate('2025-09-24T16:28:31.613Z'),
    message_count: 247
  },
  {
    _id: ObjectId('6971252104739cc1d728c128'),
    participants_count: 2,
    last_message_at: ISODate('2025-11-10T06:56:15.978Z'),
    message_count: 242
  },
  {
    _id: ObjectId('6971252104739cc1d728c130'),
    participants_count: 4,
    last_message_at: ISODate('2025-12-21T07:08:54.322Z'),
    message_count: 241
  },
  {
    _id: ObjectId('6971252104739cc1d728c138'),
    participants_count: 4,
    last_message_at: ISODate('2025-12-08T06:38:06.640Z'),
    message_count: 231
  },
  {
    _id: ObjectId('6971252104739cc1d728c121'),
    participants_count: 3,
    last_message_at: ISODate('2025-12-16T03:46:32.798Z'),
    message_count: 229
  },
  {
    _id: ObjectId('6971252104739cc1d728c10a'),
    participants_count: 2,
    last_message_at: ISODate('2025-11-04T12:12:32.267Z'),
    message_count: 227
  },
  {
    _id: ObjectId('6971252104739cc1d728c11c'),
    participants_count: 6,
    last_message_at: ISODate('2025-12-17T13:14:37.672Z'),
    message_count: 227
  },
  {
    _id: ObjectId('6971252104739cc1d728c109'),
    participants_count: 3,
    last_message_at: ISODate('2025-11-20T20:28:26.225Z'),
    message_count: 226
  }
]
```

---

## Pipeline 3 — Reaction distribution (optional)

This pipeline counts how many reactions exist for each `reaction_type`.

```js
db.conversations.aggregate([
  { $unwind: "$messages" },
  { $unwind: "$messages.reactions" },
  { $group: { _id: "$messages.reactions.reaction_type", total: { $sum: 1 } } },
  { $sort: { total: -1 } }
])
```

Output:
```js
[
  { _id: 'laugh', total: 2982 },
  { _id: 'wow', total: 2974 },
  { _id: 'angry', total: 2971 },
  { _id: 'sad', total: 2963 },
  { _id: 'like', total: 2931 },
  { _id: 'love', total: 2930 }
]
```
