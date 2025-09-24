# WASAText

This repository follows the Fantastic Coffee (Decaffeinated) template structure while preserving this project's UI and functionality.

## Project structure

* `cmd/` contains all executables (e.g., `cmd/webapi`, `cmd/healthcheck`)
* `demo/` contains a demo config file
* `doc/` contains the OpenAPI file (`doc/api.yaml`)
* `service/` contains project-specific packages (`service/api`, `service/globaltime`, `service/models`, `service/store`)
* `vendor/` contains vendored dependencies
* `webui/` contains the Vue.js frontend built with Vite; includes code for release embedding
* Top-level: `.editorconfig`, `.gitignore`, `go.mod`, `go.sum`, `LICENSE`, `README.md`, `open-node.sh`

Template reference: https://github.com/sapienzaapps/fantastic-coffee-decaffeinated

## How to build

Backend only:
```bash
go build ./cmd/webapi/
```

With embedded WebUI:
```bash
./open-node.sh
# inside the container
yarn run build-embed
exit
# outside the container
go build -tags webui ./cmd/webapi/
```

## How to run (development)

Backend only:
```bash
go run ./cmd/webapi/
```

WebUI (dev server with proxy):
```bash
./open-node.sh
# inside the container
yarn run dev
```

Equivalent npm commands (if you prefer outside the container):
```bash
(cd webui && npm ci && npm run dev)
```

## How to build for production / delivery

```bash
./open-node.sh
# inside the container
yarn run build-prod
yarn run preview
```

## Docker (HW4)

Build and run backend:
```bash
docker build -t new-wasa-backend -f Dockerfile.backend .
docker run -p 8080:8080 new-wasa-backend
```

Build and run frontend (nginx):
```bash
docker build -t new-wasa-frontend -f Dockerfile.frontend .
docker run -p 8081:80 new-wasa-frontend
```

Docker Compose (frontend + backend):
```bash
docker compose up -d --build
# Backend:  http://localhost:8080
# Frontend: http://localhost:8081
```

## CORS

CORS preflight is allowed for all origins, with `Access-Control-Max-Age: 1` second, as required.

## OpenAPI

`doc/api.yaml` contains all required operationIds: `doLogin`, `setMyUserName`, `getMyConversations`, `getConversation`, `sendMessage`, `forwardMessage`, `commentMessage`, `uncommentMessage`, `deleteMessage`, `addToGroup`, `leaveGroup`, `setGroupName`, `setMyPhoto`, `setGroupPhoto`.
