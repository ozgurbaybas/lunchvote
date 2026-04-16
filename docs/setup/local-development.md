---

## `docs/setup/local-development.md`

```md
# Local Development

## Requirements

- Go
- Docker
- Docker Compose
- PostgreSQL via Docker

## Start PostgreSQL

```bash
docker compose up -d


Apply Migrations

Run migration files in order:


docker compose exec -T postgres psql -U lunchvote -d lunchvote < migrations/000001_init_extensions.up.sql
docker compose exec -T postgres psql -U lunchvote -d lunchvote < migrations/000002_create_identity_tables.up.sql
docker compose exec -T postgres psql -U lunchvote -d lunchvote < migrations/000003_create_restaurant_tables.up.sql
docker compose exec -T postgres psql -U lunchvote -d lunchvote < migrations/000004_create_rating_tables.up.sql
docker compose exec -T postgres psql -U lunchvote -d lunchvote < migrations/000005_create_poll_tables.up.sql

Run Application
go run ./cmd/api


Verify
curl http://localhost:8080/health



