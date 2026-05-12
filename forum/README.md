# Forum

A web forum built with Go that allows users to communicate through posts and comments, org
anize content by categories, and interact via likes and dislikes.

---

## Features

- User registration and login with encrypted passwords
- Session management with cookies and expiration
- Create posts and associate them with one or more categories
- Comment on posts
- Like and dislike posts and comments
- Filter posts by category, your own posts, or posts you have liked
- All posts and comments visible to non-registered users

---

## Requirements

- Go 1.21+
- Docker
- Docker Compose

---

## Running Locally

Clone the repository:

```bash
git clone https://github.com/your-repo/forum.git
cd forum
```

Run the server:

```bash
go run cmd/main.go
```

Open your browser at:

```
http://localhost:8080
```

---

## Running with Docker

Build and start the container:

```bash
docker-compose up --build
```

Stop the container:

```bash
docker-compose down
```

The database file is persisted in `internals/db/data/forum.db` on your local machine via a
 Docker volume — your data survives container restarts.

---

## Project Structure

```
forum/
├── cmd/
│   └── main.go                  # Entry point
├── internals/
│   ├── db/
│   │   ├── db.go                # Database initialization
│   │   ├── schema.sql           # Table definitions
│   │   └── data/                # SQLite database file
│   ├── handlers/                # HTTP request handlers
│   ├── middleware/              # Authentication middleware
│   ├── models/                  # Data structs
│   ├── services/                # Business logic
│   ├── repository/              # Database queries
│   └── utils/                   # Helper functions
├── web/
│   ├── pages/                   # HTML templates
│   └── static/                  # CSS, JS, icons
├── Dockerfile
├── docker-compose.yml
└── go.mod
```

---

## Authors

- Sospeter Kinyanjui
- Juma Evans
- Silas Lelei
- Martin Odongo




