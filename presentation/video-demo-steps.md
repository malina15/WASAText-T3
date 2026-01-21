# MongoDB Video Demonstration — Step-by-Step Script

This file is used as a **personal execution script** for recording the video
demonstration required by the Database Systems project (T3).
All commands are executed from the root of the repository.

---

## 1. Start the application containers

```bash
docker compose up -d --build
docker ps
```

This step shows that the frontend and backend services start correctly using Docker.

---

## 2. Start a MongoDB container for the study

```bash
docker rm -f wasatext-mongo 2>/dev/null || true
docker run -d --name wasatext-mongo -p 27017:27017   -v "$(pwd)/mongodb-study/scripts:/scripts"   mongo:7
```

This container is used exclusively for the MongoDB schema design and benchmarking study.

---

## 3. Connect to MongoDB

```bash
docker exec -it wasatext-mongo mongosh
```

---

## 4. Select database and populate dataset

```js
use wasatext_embedded
load("/scripts/populate_embedded.js")
db.users.countDocuments()
db.conversations.countDocuments()
```

This step loads a synthetic dataset and confirms successful population.

---

## 5. Execute query BEFORE indexing

```js
db.conversations.find({})
  .sort({ last_message_at: -1 })
  .limit(20)
  .explain("executionStats")
```

Expected observation:
- Collection scan (`COLLSCAN`)
- Blocking in-memory sort
- High number of examined documents

---

## 6. Create index on last_message_at

```js
db.conversations.createIndex({ last_message_at: -1 })
```

---

## 7. Execute query AFTER indexing

```js
db.conversations.find({})
  .sort({ last_message_at: -1 })
  .limit(20)
  .explain("executionStats")
```

Expected observation:
- Index scan (`IXSCAN`)
- No blocking sort
- Reduced number of examined documents
- Improved execution time

---

## 8. Exit MongoDB shell

```js
exit
```

---

## End of video

This completes the video demonstration showing:
- Docker execution
- MongoDB dataset population
- Query benchmarking before and after indexing
