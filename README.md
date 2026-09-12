# Go User API

A small CRUD API built with Gorilla Mux, GORM v1, and MySQL.

## Configuration

The application reads its database configuration from environment variables. `DATABASE_DSN` takes precedence when set; otherwise set `DB_USER`, `DB_PASSWORD`, `DB_HOST` (defaults to `127.0.0.1:3306`), and `DB_NAME`.

```powershell
$env:DATABASE_DSN = "user:password@tcp(127.0.0.1:3306)/database?parseTime=true"
go run .
```

Set `PORT` to use a port other than `8098`.

The API expects an existing `user` table with `id`, `name`, and `email` columns. It intentionally does not run migrations automatically.

## Endpoints

| Method | Path | Body | Success response |
| --- | --- | --- | --- |
| `GET` | `/user/{id}` | — | `200` |
| `POST` | `/user/new` | `{"name":"Ada","email":"ada@example.com"}` | `201` |
| `PATCH` | `/user/update/{id}` | one or both of `name`, `email` | `200` |
| `DELETE` | `/user/delete/{id}` | — | `204` |

All error responses are JSON and include an appropriate HTTP status code. Request bodies must be a single JSON object; unknown fields and invalid email addresses are rejected.

## Sample clients

Run the server first, then execute any of these from a second terminal:

```powershell
go run ./cmd/adduser
go run ./cmd/updateuser
go run ./cmd/deleteuser
```

The update and delete examples use ID `1`; change it to an ID in your database as needed.

## Verification

```powershell
go test ./...
go vet ./...
```
