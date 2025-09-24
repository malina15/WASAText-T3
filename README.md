# WASAText

This repository contains the basic structure for Web and Software Architecture homework project. It has been described in class.

"WASAText" is a simplified messaging application for the WASA course, not suitable for a production environment. The full version can be found in more advanced repositories.

## Project structure

* `cmd/` contains all executables; Go programs here should only do "executable-stuff", like reading options from the CLI/env, etc.  
   * `cmd/healthcheck` is an example of a daemon for checking the health of servers daemons; useful when the hypervisor is not providing HTTP readiness/liveness probes (e.g., Docker engine)  
   * `cmd/webapi` contains an example of a web API server daemon
* `demo/` contains a demo config file
* `doc/` contains the documentation (usually, for APIs, this means an OpenAPI file)
* `service/` has all packages for implementing project-specific functionalities  
   * `service/api` contains an example of an API server  
   * `service/globaltime` contains a wrapper package for `time.Time` (useful in unit testing)
   * `service/models` contains data models for the chat application
   * `service/store` contains data access layer implementations
* `vendor/` is managed by Go, and contains a copy of all dependencies
* `webui/` is a web frontend in Vue.js built with Vite; it includes Go code for release embedding.

Other project files include:

* `open-node.sh` starts a new (temporary) container using `node:20` image for safe and secure web frontend development (you don't want to use `node` in your system, do you?).

## Go vendoring

This project uses Go Vendoring. You must use `go mod vendor` after changing some dependency (`go get` or `go mod tidy`) and add all files under `vendor/` directory in your commit.

For more information about vendoring:

* <https://go.dev/ref/mod#vendoring>
* <https://www.ardanlabs.com/blog/2020/04/modules-06-vendoring.html>

## Node vendoring

The `webui` uses npm. Dependencies are installed via `npm ci` (or `npm install`).

## How to set up a new project from this template

You need to:

* Change the Go module path to your module path in `go.mod`, `go.sum`, and in `*.go` files around the project
* Rewrite the API documentation `doc/api.yaml`
* If no web frontend is expected, remove `webui` and `cmd/webapi/register-webui.go`
* Update top/package comment inside `cmd/webapi/main.go` to reflect the actual project usage, goal, and general info
* Update the code in `run()` function (`cmd/webapi/main.go`) to connect to databases or external resources
* Write API code inside `service/api`, and create any further package inside `service/` (or subdirectories)

## How to build

If you're not using the WebUI, or if you don't want to embed the WebUI into the final executable, then:

```bash
go build ./cmd/webapi/
```

If you're using the WebUI and you want to embed it into the final executable:

```bash
# build frontend for embedding
(cd webui && npm ci && npm run build-embed)
# build backend with webui tag
go build -tags webui ./cmd/webapi/
```

## How to run (in development mode)

You can launch the backend only using:

```bash
go run ./cmd/webapi/
```

If you want to launch the WebUI locally with proxy to backend:

```bash
(cd webui && npm ci && npm run dev)
```

Alternatively, using the containerized Node environment:

```bash
./open-node.sh
# inside the container
npm run dev
```

## How to build for production / homework delivery

```bash
./open-node.sh
# inside the container
npm run build-prod
exit
# outside the container
go build -tags webui ./cmd/webapi/
```

To preview the production frontend locally:

```bash
(cd webui && npm run preview)
```

## Known issues

### My build works when I use `npm run dev`, however there is a Javascript crash in production/grading

Some errors in the code are somehow not shown in `vite` development mode. To preview the code that will be used in production/grading settings, use the following commands:

```bash
(cd webui && npm run build-prod)
(cd webui && npm run preview)
```

## API Endpoints

The application provides the following REST API endpoints:

### Authentication
- `POST /session` - User login
- `PUT /user/username` - Update username
- `PUT /user/photo` - Update user photo

### Conversations
- `GET /conversations` - Get user conversations
- `POST /conversations` - Create new conversation
- `GET /conversations/{id}` - Get conversation details

### Messages
- `POST /messages` - Send message
- `POST /messages/{id}/forward` - Forward message
- `POST /messages/{id}/comment` - Comment on message
- `POST /messages/{id}/uncomment` - Remove comment
- `DELETE /messages/{id}` - Delete message

### Groups
- `POST /groups/{id}/add` - Add user to group
- `POST /groups/{id}/leave` - Leave group
- `PUT /groups/{id}/name` - Update group name
- `PUT /groups/{id}/photo` - Update group photo

## License

See LICENSE.

## About

Template repository for homework project in Web and Software Architecture course @ Sapienza University of Rome - <http://gamificationlab.uniroma1.it/en/wasa/>
