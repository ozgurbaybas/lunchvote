# ADR-002: PostgreSQL Only

## Status
Accepted

## Context
The current product scope fits well into a relational model and the project explicitly requires PostgreSQL as the persistence layer.

## Decision
PostgreSQL will be the only primary database.

## Consequences
- simpler operational model
- transactional consistency across modules
- easier integration testing
- schema evolution must be managed carefully with migrations