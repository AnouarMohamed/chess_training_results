# Chess Training Results API

Go backend for tracking chess training data: players, repertoires, target lines, drill sessions, games, and reviews.

## Stack

- Go 1.23
- Gin (HTTP API)
- PostgreSQL 16
- pgx/v5
- sqlc
- Docker / Docker Compose

## Quick Start

1. Start PostgreSQL:

```bash
make up
```

2. Export required environment variables:

```bash
export DATABASE_URL="postgres://chess:chess@localhost:5432/chess?sslmode=disable"
export JWT_SECRET="change-me"
export JWT_TTL_MIN="120"
# optional (default: :8080)
export HTTP_ADDR=":8080"
```

Validation rules:
- `DATABASE_URL` is required
- `JWT_SECRET` is required
- `JWT_TTL_MIN` must be a positive integer

3. Apply migrations (if you have `migrate` CLI installed):

```bash
make migrate-up
```

4. Run API:

```bash
make run
```

API base URL: `http://localhost:8080`

## Main Endpoints

- `GET /healthz`
- `POST /auth/register`
- `POST /auth/login`
- `GET /auth/me` (requires `Authorization: Bearer <token>`)
  - Returns `userId`, `username`, and `requestId` when request tracing middleware is active

## Development Commands

```bash
make fmt
make test
make sqlc
make build
```

## SQLC

Generated files are committed under `internal/db/sqlc`.  
If SQL schema or queries change, regenerate:

```bash
sqlc generate
```
