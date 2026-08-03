# Social Network

A Facebook-like social networking platform built with:

- **Frontend:** Next.js (React)
- **Backend:** Go (net/http)
- **Database:** SQLite
- **Realtime:** Gorilla WebSocket
- **Auth:** session cookies
- **Deployment:** Docker Compose

> Note: this repository contains additional docs under `docs/` (API, architecture, ERD), plus separate frontend documentation under `frontend/`.

---

## Quick start (local)

### 1) Frontend (Next.js)

```bash
cd frontend
npm ci
npm run dev
```

Open: http://localhost:3000

### 2) Backend (Go)

```bash
cd backend
go run ./cmd/server
```

Backend will start on the configured host/port (defaults are controlled by `backend/pkg/config` and environment variables).

---

## Docker (recommended)

Start both frontend and backend containers:

```bash
docker compose up --build
```

Then:
- Frontend: http://localhost:3000
- Backend: typically http://localhost:8080

---

## Project structure

```text
backend/
  cmd/server/main.go          # HTTP + WebSocket server entrypoint
  pkg/                        # config, handlers, middleware, repositories, routes, websocket, migrations
frontend/
  app/                        # Next.js app router
  src/
    router/
    services/                # API + websocket clients
    store/                   # app state
```

---

## Backend overview

- **HTTP REST APIs** via `net/http` and a `ServeMux`.
- **WebSocket hub** for realtime features.
- **SQLite migrations** run on startup.
- **Session middleware** + **rate limiting**.

Key entrypoints:
- `backend/cmd/server/main.go`

---

## API and domain docs

- `docs/API.md` — API endpoints and request/response formats
- `docs/Architecture.md` — architecture overview
- `docs/ERD.md` — database schema / ERD

---

## Frontend docs

See:
- `frontend/README.md`
- `frontend/AGENTS.md`, `frontend/CLAUDE.md` (if present in this repo)

---

## Development notes

- Use the backend migrations (`backend/pkg/db/migrations/...`) for schema changes.
- Keep business logic in services and SQL in repositories (per architecture docs).

