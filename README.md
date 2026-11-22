# Go User API

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat-square&logo=go&logoColor=white)](https://go.dev/)
[![Gorilla Mux](https://img.shields.io/badge/Router-Gorilla%20Mux-blue?style=flat-square)](#)
[![GORM](https://img.shields.io/badge/ORM-GORM%20v1-00ADD8?style=flat-square)](#)
[![MySQL](https://img.shields.io/badge/Database-MySQL-4479A1?style=flat-square&logo=mysql&logoColor=white)](#)

A lightweight RESTful User CRUD API built with **Go**, **Gorilla Mux**, **GORM v1**, and **MySQL**, featuring strict JSON request validation, structured response envelopes, and sample Go HTTP clients.

---

## Features

- **Strict JSON Decoding:** Rejects unknown fields (`DisallowUnknownFields`), trailing/multiple JSON objects, and payloads larger than **1 MiB**.
- **Input Validation:** Verifies positive integer user IDs, rejects blank names, and validates email syntax with Go's `net/mail` parser.
- **Consistent JSON Envelopes:** Uniform `status`, `data`, and `message` responses across `2xx`, `4xx`, and `5xx` statuses.
- **Flexible DB Configuration:** Supports either a full `DATABASE_DSN` or individual `DB_*` environment variables.

---

## Prerequisites & Database Setup

- **Go** `1.21+`
- **MySQL** running locally or remotely

The API expects an existing `user` table with `id`, `name`, and `email` columns (it intentionally does not run automatic migrations):

```sql
CREATE DATABASE IF NOT EXISTS go_api;
USE go_api;

CREATE TABLE IF NOT EXISTS user (
  id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(120) NOT NULL,
  email VARCHAR(255) NOT NULL UNIQUE
);
```

---

## Configuration & Running

Database settings are loaded from environment variables (`DATABASE_DSN` takes precedence when set):

| Variable | Required | Default | Description |
| --- | :---: | --- | --- |
| `DATABASE_DSN` | Optional | — | Full MySQL DSN (overrides individual `DB_*` vars) |
| `DB_USER` | If no DSN | — | MySQL username |
| `DB_PASSWORD` | If no DSN | — | MySQL password |
| `DB_HOST` | No | `127.0.0.1:3306` | MySQL host and port |
| `DB_NAME` | If no DSN | — | Database schema name |
| `PORT` | No | `8098` | HTTP server listen port |

### PowerShell (Windows)

```powershell
$env:DATABASE_DSN = "user:password@tcp(127.0.0.1:3306)/go_api?parseTime=true"
go run .
```

### Bash (Linux / macOS)

```bash
export DATABASE_DSN="user:password@tcp(127.0.0.1:3306)/go_api?parseTime=true"
go run .
```

The server starts on `http://localhost:8098` (or `:$PORT`).

---

## API Endpoints

| Method | Path | Request Body | Success Status |
| --- | --- | --- | :---: |
| `GET` | `/user/{id}` | — | `200 OK` |
| `POST` | `/user/new` | `{"name":"Ada","email":"ada@example.com"}` | `201 Created` |
| `PATCH` | `/user/update/{id}` | One or both of `name`, `email` | `200 OK` |
| `DELETE` | `/user/delete/{id}` | — | `204 No Content` |

### Example Requests & Responses

**Create a user (`POST /user/new`):**

```bash
curl -X POST http://localhost:8098/user/new \
  -H "Content-Type: application/json" \
  -d '{"name":"Ada Lovelace","email":"ada@example.com"}'
```

```json
{
  "status": "success",
  "data": {
    "id": 1,
    "name": "Ada Lovelace",
    "email": "ada@example.com"
  },
  "message": "user created"
}
```

**Fetch a user (`GET /user/{id}`):**

```bash
curl http://localhost:8098/user/1
```

**Update a user (`PATCH /user/update/{id}`):**

```bash
curl -X PATCH http://localhost:8098/user/update/1 \
  -H "Content-Type: application/json" \
  -d '{"name":"Ada King"}'
```

**Delete a user (`DELETE /user/delete/{id}`):**

```bash
curl -X DELETE http://localhost:8098/user/delete/1
```

**Error response format (`400 Bad Request`, `404 Not Found`, `500 Internal Server Error`):**

```json
{
  "status": "error",
  "message": "user not found"
}
```

---

## Sample CLI Clients

With the server running on `:8098`, execute any of the included helper clients in another terminal:

```bash
go run ./cmd/adduser
go run ./cmd/updateuser
go run ./cmd/deleteuser
```

> **Note:** `cmd/updateuser` and `cmd/deleteuser` target user ID `1` on `http://localhost:8098`.

---

## Project Structure

```text
.
├── main.go                 # HTTP router, handlers, JSON decoding & validation
├── main_test.go            # Unit tests for request decoding and ID parsing
├── connect/
│   └── connect.go          # MySQL/GORM connection & user CRUD queries
├── structures/
│   └── structures.go       # User model & Response JSON envelope
└── cmd/
    ├── adduser/main.go     # Sample POST client
    ├── updateuser/main.go  # Sample PATCH client
    └── deleteuser/main.go  # Sample DELETE client
```

---

## Verification

```bash
go test -v ./...
go vet ./...
```
