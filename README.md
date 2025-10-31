# WASAgram

WASAgram is a simple social media sharing app with:
- authentication by identifier (session issuance)
- user profiles (nickname, follow/unfollow, ban/unban)
- photo posts (upload/delete, stream, likes, comments)
- responsive Vue 3 frontend and Go backend with SQLite

## How it works (in plain words)

- Login is by a simple identifier (no password). If the identifier doesn’t exist, the profile is created on the fly.
- Your Home feed shows posts from people you follow. If you don’t follow anyone yet, your feed will be empty until you do.
- You can like and comment on other people’s posts. You can’t like your own posts.
- Uploading a photo adds it to your Profile. Your followers will see it in their Home feed.
- You can follow or unfollow other users at any time.
- If you ban someone, they won’t appear in your feed and they won’t see your posts.
- Deleting a photo removes its likes and comments.

## Project structure

* `cmd/` contains all executables (e.g., `cmd/webapi`, `cmd/healthcheck`)
* `demo/` contains a demo config file
* `doc/` contains the OpenAPI file (`doc/api.yaml`)
* `service/` contains project-specific packages (`service/api`, `service/globaltime`, `service/models`, `service/store`)
* `webui/` contains the Vue.js frontend built with Vite; includes code for release embedding
* Top-level: `.editorconfig`, `.gitignore`, `go.mod`, `go.sum`, `LICENSE`, `README.md`, `open-node.sh`


## Quick Start with Docker Compose (Recommended)

The easiest way to run the complete application (backend + frontend):

```bash
# Clone the repository
git clone [repository-url]
cd new-wasa

# Start everything with Docker Compose
docker compose up -d --build

# Access the application:
# - Frontend: http://localhost:8081
# - Backend API: http://localhost:3000
```

To stop:
```bash
docker compose down
```

## Backend build

Backend only:
```bash
go build ./cmd/webapi/
```

The backend exposes port 3000 and reads config via flags/env (see `cmd/webapi/load-configuration.go`).

## Run in development

Backend only:
```bash
go run ./cmd/webapi/
```

WebUI (dev server):
```bash
cd webui && npm ci && npm run dev
```

To build the frontend for production locally: `cd webui && npm ci && npm run build-prod` (assets in `webui/dist`).

## Docker (HW4)

Build and run backend:
```bash
docker build -t new-wasa-backend -f Dockerfile.backend .
docker run -p 3000:3000 new-wasa-backend
```

Build and run frontend (nginx):
```bash
docker build -t new-wasa-frontend -f Dockerfile.frontend .
docker run -p 8081:80 new-wasa-frontend
```

Docker Compose (frontend + backend):
```bash
docker compose up -d --build
# Backend:  http://localhost:3000
# Frontend: http://localhost:8081
```

## Notes

- CORS is enabled with preflight support; default max-age is 1s.
- Nginx in the frontend container proxies API paths to the backend service.
- Backend listens on `0.0.0.0:3000` by default when containerized.
