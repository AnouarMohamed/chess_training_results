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
cp .env.example .env
# then export or load from your shell/env manager
export DATABASE_URL="postgres://chess:chess@localhost:5432/chess?sslmode=disable"
export JWT_SECRET="change-me"
export JWT_TTL_MIN="120"
export HTTP_ADDR=":8080"
# optional (default: dev): one of dev|test|staging|prod
export APP_ENV="dev"
```

Validation rules:
- `DATABASE_URL` is required
- `JWT_SECRET` is required
- `JWT_TTL_MIN` must be a positive integer
- `APP_ENV` must be one of `dev`, `test`, `staging`, `prod`

3. Apply migrations (if you have `migrate` CLI installed):

```bash
make migrate-up
```

4. Run API:

```bash
make run
```

API base URL: `http://localhost:8080`

`docker compose ps` will also report health status for `db` and `api` containers.

## Main Endpoints

- `GET /healthz` (returns `ok` and `env`)
- `POST /auth/register`
- `POST /auth/login`
- `GET /auth/me` (requires `Authorization: Bearer <token>`)
  - Returns `userId`, `username`, and `requestId` when request tracing middleware is active

Authorization header format for protected routes:
- `Authorization: Bearer <jwt>`

Security headers added to all responses:
- `X-Content-Type-Options: nosniff`
- `X-Frame-Options: DENY`
- `Referrer-Policy: no-referrer`

Registration input rules:
- Username: 3-32 chars, letters/numbers/underscore only
- Password: minimum 8 characters

Authentication request validation:
- Register `username`: 3-32 chars
- Register `password`: 8-128 chars
- Register `email`: optional but must be valid email format when present

## Development Commands

```bash
make fmt
make test
make test-fresh
make vet
make check
make sqlc
make build
make print-env-template
```

## SQLC

Generated files are committed under `internal/db/sqlc`.  
If SQL schema or queries change, regenerate:

```bash
sqlc generate
```
