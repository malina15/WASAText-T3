# CRUD operations (T3) — WASAText embedded schema

All commands were executed in `mongosh` on database `wasatext_embedded`.
Dataset check (as shown in the terminal):

```js
use wasatext_embedded
show collections
db.users.countDocuments()
db.conversations.countDocuments()
```

Output:
```txt
conversations
users
200
80
```

---

## Create — create conversation with first message

```js
const u1 = db.users.findOne({}, { _id: 1 })
const u2 = db.users.findOne({ _id: { $ne: u1._id } }, { _id: 1 })

const ins = db.conversations.insertOne({
  participants: [u1._id, u2._id],
  participants_count: 2,
  created_at: new Date(),
  last_message_at: new Date(),
  messages: [{
    _id: new ObjectId(),
    sender_id: u1._id,
    content: "T3 demo: hello",
    sent_at: new Date(),
    reactions: [],
    receipts: [{ user_id: u2._id, status: "delivered", at: new Date() }]
  }]
})

ins
```

Output:
```js
{
  acknowledged: true,
  insertedId: ObjectId('6971264b476e19578e284d0d')
}
```

---

## Read — Inbox view (latest conversations)

```js
db.conversations.find(
  {},
  { participants: 1, participants_count: 1, last_message_at: 1 }
).sort({ last_message_at: -1 }).limit(5)
```

Output (captured):
```js
[
  {
    _id: ObjectId('6971264b476e19578e284d0d'),
    participants: [
      ObjectId('6971252004739cc1d7284d0c'),
      ObjectId('6971252004739cc1d7284d0d')
    ],
    participants_count: 2,
    last_message_at: ISODate('2026-01-21T19:17:31.593Z')
  }
]
```

---

## Update — send a new message (append) + update last_message_at

```js
const convId = ins.insertedId

db.conversations.updateOne(
  { _id: convId },
  {
    $push: {
      messages: {
        _id: new ObjectId(),
        sender_id: u1._id,
        content: "T3 demo: second message",
        sent_at: new Date(),
        reactions: [],
        receipts: []
      }
    },
    $set: { last_message_at: new Date() }
  }
)
```

Output:
```js
{
  acknowledged: true,
  insertedId: null,
  matchedCount: 1,
  modifiedCount: 1,
  upsertedCount: 0
}
```

Verification (last 2 messages):
```js
db.conversations.findOne(
  { _id: convId },
  { messages: { $slice: -2 }, last_message_at: 1 }
)
```

---

## Update — react to the latest message (optional but realistic)

```js
const lastMsgId = db.conversations.findOne(
  { _id: convId },
  { messages: { $slice: -1 } }
).messages[0]._id

db.conversations.updateOne(
  { _id: convId, "messages._id": lastMsgId },
  { $push: { "messages.$.reactions": { user_id: u2._id, type: "love", at: new Date() } } }
)
```

---

## Delete — delete the demo conversation (cleanup)

```js
db.conversations.deleteOne({ _id: convId })
```

