---

## `docs/setup/testing.md`

```md
# Testing Strategy

LunchVote uses multiple test layers.

## 1. Domain Unit Tests

Purpose:
- verify business rules
- validate entities and aggregates

## 2. Application Tests

Purpose:
- verify use case orchestration
- test command/query behavior
- validate domain + repository interaction with in-memory doubles

## 3. Repository Integration Tests

Purpose:
- verify PostgreSQL repository implementations
- validate persistence mappings
- validate constraints and retrieval logic

## 4. HTTP Handler Tests

Purpose:
- verify request parsing
- verify response status mapping
- verify endpoint behavior

## Run All Tests

```bash
go test ./...


Run Module-Specific Tests

go test ./modules/identity/...
go test ./modules/restaurant/...
go test ./modules/rating/...
go test ./modules/poll/...
go test ./modules/recommendation/...