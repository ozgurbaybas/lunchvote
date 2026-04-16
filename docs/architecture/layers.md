# Layered Structure

Each business module follows the same layering model.

## domain

Contains:
- entities
- aggregates
- value objects
- business rules
- repository contracts

Must remain framework independent.

## application

Contains:
- use cases
- commands / queries
- orchestration logic
- cross-module coordination

## infrastructure

Contains:
- PostgreSQL repository implementations
- persistence details

## interfaces/http

Contains:
- handlers
- request/response DTOs
- HTTP status mapping
- basic request validation