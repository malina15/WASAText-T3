# MongoDB Video Demonstration — Step-by-Step Script

This file is used as a **personal execution script** for recording the video
demonstration required by the Database Systems project (T3).
All commands are executed from the repository root, using the VS Code integrated terminal,
unless explicitly stated otherwise.

---

## 1. Start the application using Docker Compose

```bash
docker compose up -d --build
docker ps
```

At this point, I visually confirm that:
- the frontend container is running on port 8081
- the backend container is running on port 3000

This shows that the application starts correctly in a containerized environment.

---

## 2. Start a dedicated MongoDB container for the study

```bash
docker rm -f wasatext-mongo 2>/dev/null || true
docker run -d --name wasatext-mongo -p 27017:27017   -v "$(pwd)/mongodb-study/scripts:/scripts"   mongo:7
```

This MongoDB instance is used **only** for the MongoDB schema design,
indexing, and benchmarking experiments required by T3.

---

## 3. Connect to MongoDB using mongosh

```bash
docker exec -it wasatext-mongo mongosh
```

This opens an interactive MongoDB shell inside the container.

---

## 4. Select database and populate the dataset

```js
use wasatext_embedded
load("/scripts/populate_embedded.js")
db.users.countDocuments()
db.conversations.countDocuments()
```

Here I:
- switch to the `wasatext_embedded` database
- load the synthetic dataset generation script
- verify that users and conversations were successfully inserted

---

## 5. Execute query BEFORE creating the index

```js
db.conversations.find({})
  .sort({ last_message_at: -1 })
  .limit(20)
  .explain("executionStats")
```

This query simulates the **Inbox view**, fetching the most recent conversations.
At this stage, the query runs **without any index**, which is visible in the
execution plan.

Expected behavior:
- collection scan
- in-memory sort
- higher execution cost

---

## 6. Create index on `last_message_at`

```js
db.conversations.createIndex({ last_message_at: -1 })
```

This index matches the access pattern used by the Inbox query.

---

## 7. Execute the same query AFTER indexing

```js
db.conversations.find({})
  .sort({ last_message_at: -1 })
  .limit(20)
  .explain("executionStats")
```

Now the execution plan shows:
- index scan on `last_message_at`
- no blocking sort
- fewer examined documents
- improved performance

This demonstrates the impact of indexing on a realistic workload.

---

## 8. Exit the MongoDB shell

```js
exit
```

---

## End of recording

The video demonstrates:
- Docker-based execution of the application
- MongoDB dataset population
- Query execution before and after indexing
- Measurable performance improvement through indexing

