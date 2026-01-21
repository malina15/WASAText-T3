# Replication (T3) — MongoDB Replica Set Demo

This section demonstrates MongoDB replication using a 3-node replica set deployed with Docker.
The goal is to show leader election and automatic failover (high availability).

---

## Setup (Docker)

Commands were executed in the VS Code terminal from the repository root.

```bash
docker network create mongo-rs

docker run -d --name mongo1 --net mongo-rs -p 27117:27017 mongo:7 mongod --replSet rs0 --bind_ip_all
docker run -d --name mongo2 --net mongo-rs -p 27118:27017 mongo:7 mongod --replSet rs0 --bind_ip_all
docker run -d --name mongo3 --net mongo-rs -p 27119:27017 mongo:7 mongod --replSet rs0 --bind_ip_all
```

---

## Replica set initialization (mongosh on mongo1)

```js
rs.initiate({
  _id: "rs0",
  members: [
    { _id: 0, host: "mongo1:27017" },
    { _id: 1, host: "mongo2:27017" },
    { _id: 2, host: "mongo3:27017" }
  ]
})
```

Initial status after election:

```js
rs.status().members.map(m => ({ name: m.name, stateStr: m.stateStr, health: m.health }))
```

Expected result:
- one PRIMARY node
- two SECONDARY nodes

---

## Failover test

The PRIMARY node was stopped:

```bash
docker stop mongo1
```

Status after stopping the PRIMARY:

```bash
docker exec -it mongo2 mongosh --quiet --eval 'rs.status().members.map(m => ({name:m.name,stateStr:m.stateStr}))'
```

Output captured during the demo:

```js
[
  { name: 'mongo1:27017', stateStr: '(not reachable/healthy)' },
  { name: 'mongo2:27017', stateStr: 'PRIMARY' },
  { name: 'mongo3:27017', stateStr: 'SECONDARY' }
]
```

---

## Conclusion

When the PRIMARY node becomes unavailable, MongoDB automatically elects a new PRIMARY.
This demonstrates high availability and fault tolerance provided by replication.
